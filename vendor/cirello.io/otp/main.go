// Command otp manages one-time passwords tokens, protecting them with a local
// private key (usually $HOME/.ssh/id_rsa) and storing its information in a
// encrypted db (usually at $HOME/.ssh/auth.db).
package main // import "cirello.io/otp"

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
	"strings"
	"text/tabwriter"
	"time"

	otp "github.com/hgfischer/go-otp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli"
	"rsc.io/qr"
)

var homeDir string

func init() {
	log.SetPrefix("")
	log.SetFlags(0)

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	homeDir = usr.HomeDir
}

func main() {
	if os.Getenv("OTP_TRACE") != "" {
		w, err := os.Create("otp.trace")
		if err != nil {
			log.Fatal("cannot create application tracing:", err)
		}
		err = trace.Start(w)
		if err != nil {
			log.Fatal("cannot start application tracing:", err)
		}
		defer trace.Stop()

		fcpu, err := os.Create("otp.cpuprofile")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(fcpu); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	app := cli.NewApp()
	app.Name = "OTP client"
	app.Usage = "command interface"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "db",
			Value:  filepath.Join(homeDir, ".ssh", "auth.db"),
			EnvVar: "OTP_DB",
		},
		cli.StringFlag{
			Name:   "private-key",
			Value:  filepath.Join(homeDir, ".ssh", "id_rsa"),
			EnvVar: "OTP_PRIVKEY",
		},
	}
	app.Commands = []cli.Command{
		initdb(),
		add(),
		get(),
		list(),
		genqr(),
		rm(),
		servehttp(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("error: %v", err)
	}

	if os.Getenv("OTP_TRACE") != "" {
		fmem, err := os.Create("otp.memprofile")
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(fmem); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		fmem.Close()
	}
}

func initdb() cli.Command {
	return cli.Command{
		Name:  "init",
		Usage: "initialize the OTP database",
		Action: func(c *cli.Context) error {
			db, err := sql.Open("sqlite3", c.GlobalString("db"))
			if err != nil {
				return err
			}
			defer db.Close()

			queries := []string{
				"CREATE TABLE IF NOT EXISTS `otps` (`id` INTEGER PRIMARY KEY, `account` char, `issuer` char, `password` blob);",
				"CREATE UNIQUE INDEX `otps_account_issuer` ON `otps`(`account`, `issuer`);",
			}

			for _, q := range queries {
				_, err = db.Exec(q)
				if err != nil {
					return err
				}
			}

			log.Println("database initialized")
			return nil
		},
	}
}

func add() cli.Command {
	return cli.Command{
		Name:      "add",
		Usage:     "a new OTP key",
		ArgsUsage: "`secret` `issuer` `account-name`",
		Action: func(c *cli.Context) error {
			priv, err := privkeyfile(c.GlobalString("private-key"))
			if err != nil {
				return err
			}

			secretkey := c.Args().Get(0)
			issuer := c.Args().Get(1)
			account := c.Args().Get(2)

			switch {
			case secretkey == "":
				return errors.New("secret key is missing")
			case issuer == "":
				return errors.New("issuer is missing")
			case account == "":
				return errors.New("account name is missing")
			}

			enckey, err := priv.encrypted([]byte(secretkey), cryptlabel(account, issuer))
			if err != nil {
				return err
			}

			db, err := sql.Open("sqlite3", c.GlobalString("db"))
			if err != nil {
				return err
			}
			defer db.Close()

			_, err = db.Exec("REPLACE INTO `otps` (`issuer`, `account`, `password`) VALUES (?, ?, ?);", issuer, account, enckey)
			return err
		},
	}
}

func get() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "generate OTP",
		Action: func(c *cli.Context) error {
			return load(c, os.Stdout)
		},
	}
}

func servehttp() cli.Command {
	return cli.Command{
		Name:  "http",
		Usage: "serve OTP in a HTTP interface",
		Action: func(c *cli.Context) error {
			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "<html><body><pre>")
				load(c, w)
				fmt.Fprintln(w, "</pre></body></html>")
			})
			http.ListenAndServe(":9999", nil)
			return nil
		},
	}
}

func load(c *cli.Context, w io.Writer) error {
	priv, err := privkeyfile(c.GlobalString("private-key"))
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", c.GlobalString("db"))
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query("SELECT `account`, `issuer`, `password` FROM `otps` ORDER BY `account` ASC, `issuer` ASC;")
	if err != nil {
		return err
	}
	defer rows.Close()

	tabw := tabwriter.NewWriter(w, 8, 8, 2, ' ', 0)
	defer tabw.Flush()
	fmt.Fprintln(tabw, "account\tissuer\texpiration\tcode")

	for rows.Next() {
		var account, issuer string
		var pw []byte
		rows.Scan(&account, &issuer, &pw)

		decrypted, err := priv.decrypted(pw, cryptlabel(account, issuer))
		if err != nil {
			return err
		}

		key := strings.ToUpper(strings.Replace(string(decrypted), " ", "", -1))
		totp := &otp.TOTP{Secret: key, IsBase32Secret: true}
		token := totp.Get()

		fmt.Fprintln(
			tabw,
			fmt.Sprintf("%s\t%s\t%vs\t%s",
				account,
				issuer,
				(30-time.Now().Unix()%30),
				token),
		)
	}

	return nil
}

func list() cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list all keys",
		Action: func(c *cli.Context) error {
			db, err := sql.Open("sqlite3", c.GlobalString("db"))
			if err != nil {
				return err
			}
			defer db.Close()

			rows, err := db.Query("SELECT account, issuer FROM `otps` ORDER BY account ASC, issuer ASC;")
			if err != nil {
				return err
			}
			defer rows.Close()

			w := tabwriter.NewWriter(os.Stdout, 8, 8, 2, ' ', 0)
			defer w.Flush()
			fmt.Fprintln(w, "account\tissuer")

			for rows.Next() {
				var account, issuer string
				rows.Scan(&account, &issuer)
				fmt.Fprintln(w, fmt.Sprintf("%s\t%s", account, issuer))
			}

			return nil
		},
	}
}

func genqr() cli.Command {
	return cli.Command{
		Name:  "qr",
		Usage: "generate QR codes",
		Action: func(c *cli.Context) error {
			priv, err := privkeyfile(c.GlobalString("private-key"))
			if err != nil {
				return err
			}

			db, err := sql.Open("sqlite3", c.GlobalString("db"))
			if err != nil {
				return err
			}
			defer db.Close()

			rows, err := db.Query("SELECT `account`, `issuer`, `password` FROM `otps` ORDER BY `account` ASC, `issuer` ASC;")
			if err != nil {
				return err
			}
			defer rows.Close()

			w := tabwriter.NewWriter(os.Stdout, 8, 8, 2, ' ', 0)
			defer w.Flush()
			fmt.Fprintln(w, "account\tissuer\tfile")

			for rows.Next() {
				var account, issuer string
				var pw []byte
				rows.Scan(&account, &issuer, &pw)

				decrypted, err := priv.decrypted(pw, cryptlabel(account, issuer))
				if err != nil {
					return err
				}

				qrfn, err := generateQR(issuer, account, string(decrypted))
				if err != nil {
					fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s", account, issuer, err))
					continue
				}
				fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s", account, issuer, qrfn))
			}

			return nil
		},
	}
}

func rm() cli.Command {
	return cli.Command{
		Name:      "rm",
		Usage:     "delete a OTP key",
		ArgsUsage: "`issuer` `account-name`",
		Action: func(c *cli.Context) error {
			issuer := c.Args().Get(0)
			account := c.Args().Get(1)

			switch {
			case issuer == "":
				return errors.New("issuer is missing")
			case account == "":
				return errors.New("account name is missing")
			}

			db, err := sql.Open("sqlite3", c.GlobalString("db"))
			if err != nil {
				return err
			}
			defer db.Close()

			_, err = db.Exec("DELETE FROM `otps` WHERE `issuer` = ? AND `account` = ?;", issuer, account)
			return err
		},
	}
}

type privkey struct {
	*rsa.PrivateKey
}

func privkeyfile(fn string) (*privkey, error) {
	pemdata, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, fmt.Errorf("cannot read key file: %s", err)
	}

	block, _ := pem.Decode(pemdata)
	if block == nil {
		return nil, errors.New("key data is not PEM encoded")
	}

	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		return nil, fmt.Errorf("mismatched key type. got: %q want: %q", got, want)
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %s", err)
	}

	return &privkey{PrivateKey: priv}, nil
}

func (p privkey) encrypted(in, label []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, &p.PublicKey, in, label)
}

func (p privkey) decrypted(in, label []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, p.PrivateKey, in, label)
}

func cryptlabel(account, issuer string) []byte {
	return []byte(fmt.Sprint(account, issuer))
}

func generateQR(issuer, account, password string) (string, error) {
	otpauth := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", issuer, account, password, issuer)
	code, err := qr.Encode(otpauth, qr.H)
	if err != nil {
		return "", err
	}

	img, _, err := image.Decode(bytes.NewReader(code.PNG()))
	if err != nil {
		panic(err)
	}

	fn := fmt.Sprintf("otp-qr-%s-%s.png", issuer, account)
	out, err := os.Create(fn)
	if err != nil {
		return "", err
	}

	err = png.Encode(out, img)
	if err != nil {
		return "", err
	}

	return fn, nil
}
