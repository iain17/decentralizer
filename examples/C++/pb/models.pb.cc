// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: models.proto

#define INTERNAL_SUPPRESS_PROTOBUF_FIELD_DEPRECATION
#include "models.pb.h"

#include <algorithm>

#include <google/protobuf/stubs/common.h>
#include <google/protobuf/stubs/port.h>
#include <google/protobuf/stubs/once.h>
#include <google/protobuf/io/coded_stream.h>
#include <google/protobuf/wire_format_lite_inl.h>
#include <google/protobuf/descriptor.h>
#include <google/protobuf/generated_message_reflection.h>
#include <google/protobuf/reflection_ops.h>
#include <google/protobuf/wire_format.h>
// @@protoc_insertion_point(includes)

namespace pb {

namespace {

const ::google::protobuf::Descriptor* peer_descriptor_ = NULL;
const ::google::protobuf::internal::GeneratedMessageReflection*
  peer_reflection_ = NULL;
const ::google::protobuf::Descriptor* peer_DetailsEntry_descriptor_ = NULL;

}  // namespace


void protobuf_AssignDesc_models_2eproto() GOOGLE_ATTRIBUTE_COLD;
void protobuf_AssignDesc_models_2eproto() {
  protobuf_AddDesc_models_2eproto();
  const ::google::protobuf::FileDescriptor* file =
    ::google::protobuf::DescriptorPool::generated_pool()->FindFileByName(
      "models.proto");
  GOOGLE_CHECK(file != NULL);
  peer_descriptor_ = file->message_type(0);
  static const int peer_offsets_[4] = {
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(peer, ip_),
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(peer, port_),
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(peer, rpcport_),
    GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(peer, details_),
  };
  peer_reflection_ =
    ::google::protobuf::internal::GeneratedMessageReflection::NewGeneratedMessageReflection(
      peer_descriptor_,
      peer::internal_default_instance(),
      peer_offsets_,
      -1,
      -1,
      -1,
      sizeof(peer),
      GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(peer, _internal_metadata_));
  peer_DetailsEntry_descriptor_ = peer_descriptor_->nested_type(0);
}

namespace {

GOOGLE_PROTOBUF_DECLARE_ONCE(protobuf_AssignDescriptors_once_);
void protobuf_AssignDescriptorsOnce() {
  ::google::protobuf::GoogleOnceInit(&protobuf_AssignDescriptors_once_,
                 &protobuf_AssignDesc_models_2eproto);
}

void protobuf_RegisterTypes(const ::std::string&) GOOGLE_ATTRIBUTE_COLD;
void protobuf_RegisterTypes(const ::std::string&) {
  protobuf_AssignDescriptorsOnce();
  ::google::protobuf::MessageFactory::InternalRegisterGeneratedMessage(
      peer_descriptor_, peer::internal_default_instance());
  ::google::protobuf::MessageFactory::InternalRegisterGeneratedMessage(
        peer_DetailsEntry_descriptor_,
        ::google::protobuf::internal::MapEntry<
            ::std::string,
            ::std::string,
            ::google::protobuf::internal::WireFormatLite::TYPE_STRING,
            ::google::protobuf::internal::WireFormatLite::TYPE_STRING,
            0>::CreateDefaultInstance(
                peer_DetailsEntry_descriptor_));
}

}  // namespace

void protobuf_ShutdownFile_models_2eproto() {
  peer_default_instance_.Shutdown();
  delete peer_reflection_;
}

void protobuf_InitDefaults_models_2eproto_impl() {
  GOOGLE_PROTOBUF_VERIFY_VERSION;

  ::google::protobuf::internal::GetEmptyString();
  peer_default_instance_.DefaultConstruct();
  ::google::protobuf::internal::GetEmptyString();
  peer_default_instance_.get_mutable()->InitAsDefaultInstance();
}

GOOGLE_PROTOBUF_DECLARE_ONCE(protobuf_InitDefaults_models_2eproto_once_);
void protobuf_InitDefaults_models_2eproto() {
  ::google::protobuf::GoogleOnceInit(&protobuf_InitDefaults_models_2eproto_once_,
                 &protobuf_InitDefaults_models_2eproto_impl);
}
void protobuf_AddDesc_models_2eproto_impl() {
  GOOGLE_PROTOBUF_VERIFY_VERSION;

  protobuf_InitDefaults_models_2eproto();
  ::google::protobuf::DescriptorPool::InternalAddGeneratedFile(
    "\n\014models.proto\022\002pb\"\211\001\n\004peer\022\n\n\002ip\030\001 \001(\t\022"
    "\014\n\004port\030\002 \001(\r\022\017\n\007rpcPort\030\003 \001(\r\022&\n\007detail"
    "s\030\004 \003(\0132\025.pb.peer.DetailsEntry\032.\n\014Detail"
    "sEntry\022\013\n\003key\030\001 \001(\t\022\r\n\005value\030\002 \001(\t:\0028\001B\002"
    "H\001b\006proto3", 170);
  ::google::protobuf::MessageFactory::InternalRegisterGeneratedFile(
    "models.proto", &protobuf_RegisterTypes);
  ::google::protobuf::internal::OnShutdown(&protobuf_ShutdownFile_models_2eproto);
}

GOOGLE_PROTOBUF_DECLARE_ONCE(protobuf_AddDesc_models_2eproto_once_);
void protobuf_AddDesc_models_2eproto() {
  ::google::protobuf::GoogleOnceInit(&protobuf_AddDesc_models_2eproto_once_,
                 &protobuf_AddDesc_models_2eproto_impl);
}
// Force AddDescriptors() to be called at static initialization time.
struct StaticDescriptorInitializer_models_2eproto {
  StaticDescriptorInitializer_models_2eproto() {
    protobuf_AddDesc_models_2eproto();
  }
} static_descriptor_initializer_models_2eproto_;

namespace {

static void MergeFromFail(int line) GOOGLE_ATTRIBUTE_COLD GOOGLE_ATTRIBUTE_NORETURN;
static void MergeFromFail(int line) {
  ::google::protobuf::internal::MergeFromFail(__FILE__, line);
}

}  // namespace


// ===================================================================

#if !defined(_MSC_VER) || _MSC_VER >= 1900
const int peer::kIpFieldNumber;
const int peer::kPortFieldNumber;
const int peer::kRpcPortFieldNumber;
const int peer::kDetailsFieldNumber;
#endif  // !defined(_MSC_VER) || _MSC_VER >= 1900

peer::peer()
  : ::google::protobuf::Message(), _internal_metadata_(NULL) {
  if (this != internal_default_instance()) protobuf_InitDefaults_models_2eproto();
  SharedCtor();
  // @@protoc_insertion_point(constructor:pb.peer)
}

void peer::InitAsDefaultInstance() {
}

peer::peer(const peer& from)
  : ::google::protobuf::Message(),
    _internal_metadata_(NULL) {
  SharedCtor();
  UnsafeMergeFrom(from);
  // @@protoc_insertion_point(copy_constructor:pb.peer)
}

void peer::SharedCtor() {
  details_.SetAssignDescriptorCallback(
      protobuf_AssignDescriptorsOnce);
  details_.SetEntryDescriptor(
      &::pb::peer_DetailsEntry_descriptor_);
  ip_.UnsafeSetDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
  ::memset(&port_, 0, reinterpret_cast<char*>(&rpcport_) -
    reinterpret_cast<char*>(&port_) + sizeof(rpcport_));
  _cached_size_ = 0;
}

peer::~peer() {
  // @@protoc_insertion_point(destructor:pb.peer)
  SharedDtor();
}

void peer::SharedDtor() {
  ip_.DestroyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}

void peer::SetCachedSize(int size) const {
  GOOGLE_SAFE_CONCURRENT_WRITES_BEGIN();
  _cached_size_ = size;
  GOOGLE_SAFE_CONCURRENT_WRITES_END();
}
const ::google::protobuf::Descriptor* peer::descriptor() {
  protobuf_AssignDescriptorsOnce();
  return peer_descriptor_;
}

const peer& peer::default_instance() {
  protobuf_InitDefaults_models_2eproto();
  return *internal_default_instance();
}

::google::protobuf::internal::ExplicitlyConstructed<peer> peer_default_instance_;

peer* peer::New(::google::protobuf::Arena* arena) const {
  peer* n = new peer;
  if (arena != NULL) {
    arena->Own(n);
  }
  return n;
}

void peer::Clear() {
// @@protoc_insertion_point(message_clear_start:pb.peer)
#if defined(__clang__)
#define ZR_HELPER_(f) \
  _Pragma("clang diagnostic push") \
  _Pragma("clang diagnostic ignored \"-Winvalid-offsetof\"") \
  __builtin_offsetof(peer, f) \
  _Pragma("clang diagnostic pop")
#else
#define ZR_HELPER_(f) reinterpret_cast<char*>(\
  &reinterpret_cast<peer*>(16)->f)
#endif

#define ZR_(first, last) do {\
  ::memset(&(first), 0,\
           ZR_HELPER_(last) - ZR_HELPER_(first) + sizeof(last));\
} while (0)

  ZR_(port_, rpcport_);
  ip_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());

#undef ZR_HELPER_
#undef ZR_

  details_.Clear();
}

bool peer::MergePartialFromCodedStream(
    ::google::protobuf::io::CodedInputStream* input) {
#define DO_(EXPRESSION) if (!GOOGLE_PREDICT_TRUE(EXPRESSION)) goto failure
  ::google::protobuf::uint32 tag;
  // @@protoc_insertion_point(parse_start:pb.peer)
  for (;;) {
    ::std::pair< ::google::protobuf::uint32, bool> p = input->ReadTagWithCutoff(127);
    tag = p.first;
    if (!p.second) goto handle_unusual;
    switch (::google::protobuf::internal::WireFormatLite::GetTagFieldNumber(tag)) {
      // optional string ip = 1;
      case 1: {
        if (tag == 10) {
          DO_(::google::protobuf::internal::WireFormatLite::ReadString(
                input, this->mutable_ip()));
          DO_(::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
            this->ip().data(), this->ip().length(),
            ::google::protobuf::internal::WireFormatLite::PARSE,
            "pb.peer.ip"));
        } else {
          goto handle_unusual;
        }
        if (input->ExpectTag(16)) goto parse_port;
        break;
      }

      // optional uint32 port = 2;
      case 2: {
        if (tag == 16) {
         parse_port:

          DO_((::google::protobuf::internal::WireFormatLite::ReadPrimitive<
                   ::google::protobuf::uint32, ::google::protobuf::internal::WireFormatLite::TYPE_UINT32>(
                 input, &port_)));
        } else {
          goto handle_unusual;
        }
        if (input->ExpectTag(24)) goto parse_rpcPort;
        break;
      }

      // optional uint32 rpcPort = 3;
      case 3: {
        if (tag == 24) {
         parse_rpcPort:

          DO_((::google::protobuf::internal::WireFormatLite::ReadPrimitive<
                   ::google::protobuf::uint32, ::google::protobuf::internal::WireFormatLite::TYPE_UINT32>(
                 input, &rpcport_)));
        } else {
          goto handle_unusual;
        }
        if (input->ExpectTag(34)) goto parse_details;
        break;
      }

      // map<string, string> details = 4;
      case 4: {
        if (tag == 34) {
         parse_details:
          DO_(input->IncrementRecursionDepth());
         parse_loop_details:
          peer_DetailsEntry::Parser< ::google::protobuf::internal::MapField<
              ::std::string, ::std::string,
              ::google::protobuf::internal::WireFormatLite::TYPE_STRING,
              ::google::protobuf::internal::WireFormatLite::TYPE_STRING,
              0 >,
            ::google::protobuf::Map< ::std::string, ::std::string > > parser(&details_);
          DO_(::google::protobuf::internal::WireFormatLite::ReadMessageNoVirtual(
              input, &parser));
          DO_(::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
            parser.key().data(), parser.key().length(),
            ::google::protobuf::internal::WireFormatLite::PARSE,
            "pb.peer.DetailsEntry.key"));
          DO_(::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
            parser.value().data(), parser.value().length(),
            ::google::protobuf::internal::WireFormatLite::PARSE,
            "pb.peer.DetailsEntry.value"));
        } else {
          goto handle_unusual;
        }
        if (input->ExpectTag(34)) goto parse_loop_details;
        input->UnsafeDecrementRecursionDepth();
        if (input->ExpectAtEnd()) goto success;
        break;
      }

      default: {
      handle_unusual:
        if (tag == 0 ||
            ::google::protobuf::internal::WireFormatLite::GetTagWireType(tag) ==
            ::google::protobuf::internal::WireFormatLite::WIRETYPE_END_GROUP) {
          goto success;
        }
        DO_(::google::protobuf::internal::WireFormatLite::SkipField(input, tag));
        break;
      }
    }
  }
success:
  // @@protoc_insertion_point(parse_success:pb.peer)
  return true;
failure:
  // @@protoc_insertion_point(parse_failure:pb.peer)
  return false;
#undef DO_
}

void peer::SerializeWithCachedSizes(
    ::google::protobuf::io::CodedOutputStream* output) const {
  // @@protoc_insertion_point(serialize_start:pb.peer)
  // optional string ip = 1;
  if (this->ip().size() > 0) {
    ::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
      this->ip().data(), this->ip().length(),
      ::google::protobuf::internal::WireFormatLite::SERIALIZE,
      "pb.peer.ip");
    ::google::protobuf::internal::WireFormatLite::WriteStringMaybeAliased(
      1, this->ip(), output);
  }

  // optional uint32 port = 2;
  if (this->port() != 0) {
    ::google::protobuf::internal::WireFormatLite::WriteUInt32(2, this->port(), output);
  }

  // optional uint32 rpcPort = 3;
  if (this->rpcport() != 0) {
    ::google::protobuf::internal::WireFormatLite::WriteUInt32(3, this->rpcport(), output);
  }

  // map<string, string> details = 4;
  if (!this->details().empty()) {
    typedef ::google::protobuf::Map< ::std::string, ::std::string >::const_pointer
        ConstPtr;
    typedef ConstPtr SortItem;
    typedef ::google::protobuf::internal::CompareByDerefFirst<SortItem> Less;
    struct Utf8Check {
      static void Check(ConstPtr p) {
        ::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
          p->first.data(), p->first.length(),
          ::google::protobuf::internal::WireFormatLite::SERIALIZE,
          "pb.peer.DetailsEntry.key");
        ::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
          p->second.data(), p->second.length(),
          ::google::protobuf::internal::WireFormatLite::SERIALIZE,
          "pb.peer.DetailsEntry.value");
      }
    };

    if (output->IsSerializationDeterminstic() &&
        this->details().size() > 1) {
      ::google::protobuf::scoped_array<SortItem> items(
          new SortItem[this->details().size()]);
      typedef ::google::protobuf::Map< ::std::string, ::std::string >::size_type size_type;
      size_type n = 0;
      for (::google::protobuf::Map< ::std::string, ::std::string >::const_iterator
          it = this->details().begin();
          it != this->details().end(); ++it, ++n) {
        items[n] = SortItem(&*it);
      }
      ::std::sort(&items[0], &items[n], Less());
      ::google::protobuf::scoped_ptr<peer_DetailsEntry> entry;
      for (size_type i = 0; i < n; i++) {
        entry.reset(details_.NewEntryWrapper(
            items[i]->first, items[i]->second));
        ::google::protobuf::internal::WireFormatLite::WriteMessageMaybeToArray(
            4, *entry, output);
        Utf8Check::Check(items[i]);
      }
    } else {
      ::google::protobuf::scoped_ptr<peer_DetailsEntry> entry;
      for (::google::protobuf::Map< ::std::string, ::std::string >::const_iterator
          it = this->details().begin();
          it != this->details().end(); ++it) {
        entry.reset(details_.NewEntryWrapper(
            it->first, it->second));
        ::google::protobuf::internal::WireFormatLite::WriteMessageMaybeToArray(
            4, *entry, output);
        Utf8Check::Check(&*it);
      }
    }
  }

  // @@protoc_insertion_point(serialize_end:pb.peer)
}

::google::protobuf::uint8* peer::InternalSerializeWithCachedSizesToArray(
    bool deterministic, ::google::protobuf::uint8* target) const {
  (void)deterministic; // Unused
  // @@protoc_insertion_point(serialize_to_array_start:pb.peer)
  // optional string ip = 1;
  if (this->ip().size() > 0) {
    ::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
      this->ip().data(), this->ip().length(),
      ::google::protobuf::internal::WireFormatLite::SERIALIZE,
      "pb.peer.ip");
    target =
      ::google::protobuf::internal::WireFormatLite::WriteStringToArray(
        1, this->ip(), target);
  }

  // optional uint32 port = 2;
  if (this->port() != 0) {
    target = ::google::protobuf::internal::WireFormatLite::WriteUInt32ToArray(2, this->port(), target);
  }

  // optional uint32 rpcPort = 3;
  if (this->rpcport() != 0) {
    target = ::google::protobuf::internal::WireFormatLite::WriteUInt32ToArray(3, this->rpcport(), target);
  }

  // map<string, string> details = 4;
  if (!this->details().empty()) {
    typedef ::google::protobuf::Map< ::std::string, ::std::string >::const_pointer
        ConstPtr;
    typedef ConstPtr SortItem;
    typedef ::google::protobuf::internal::CompareByDerefFirst<SortItem> Less;
    struct Utf8Check {
      static void Check(ConstPtr p) {
        ::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
          p->first.data(), p->first.length(),
          ::google::protobuf::internal::WireFormatLite::SERIALIZE,
          "pb.peer.DetailsEntry.key");
        ::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
          p->second.data(), p->second.length(),
          ::google::protobuf::internal::WireFormatLite::SERIALIZE,
          "pb.peer.DetailsEntry.value");
      }
    };

    if (deterministic &&
        this->details().size() > 1) {
      ::google::protobuf::scoped_array<SortItem> items(
          new SortItem[this->details().size()]);
      typedef ::google::protobuf::Map< ::std::string, ::std::string >::size_type size_type;
      size_type n = 0;
      for (::google::protobuf::Map< ::std::string, ::std::string >::const_iterator
          it = this->details().begin();
          it != this->details().end(); ++it, ++n) {
        items[n] = SortItem(&*it);
      }
      ::std::sort(&items[0], &items[n], Less());
      ::google::protobuf::scoped_ptr<peer_DetailsEntry> entry;
      for (size_type i = 0; i < n; i++) {
        entry.reset(details_.NewEntryWrapper(
            items[i]->first, items[i]->second));
        target = ::google::protobuf::internal::WireFormatLite::
                   InternalWriteMessageNoVirtualToArray(
                       4, *entry, deterministic, target);
;
        Utf8Check::Check(items[i]);
      }
    } else {
      ::google::protobuf::scoped_ptr<peer_DetailsEntry> entry;
      for (::google::protobuf::Map< ::std::string, ::std::string >::const_iterator
          it = this->details().begin();
          it != this->details().end(); ++it) {
        entry.reset(details_.NewEntryWrapper(
            it->first, it->second));
        target = ::google::protobuf::internal::WireFormatLite::
                   InternalWriteMessageNoVirtualToArray(
                       4, *entry, deterministic, target);
;
        Utf8Check::Check(&*it);
      }
    }
  }

  // @@protoc_insertion_point(serialize_to_array_end:pb.peer)
  return target;
}

size_t peer::ByteSizeLong() const {
// @@protoc_insertion_point(message_byte_size_start:pb.peer)
  size_t total_size = 0;

  // optional string ip = 1;
  if (this->ip().size() > 0) {
    total_size += 1 +
      ::google::protobuf::internal::WireFormatLite::StringSize(
        this->ip());
  }

  // optional uint32 port = 2;
  if (this->port() != 0) {
    total_size += 1 +
      ::google::protobuf::internal::WireFormatLite::UInt32Size(
        this->port());
  }

  // optional uint32 rpcPort = 3;
  if (this->rpcport() != 0) {
    total_size += 1 +
      ::google::protobuf::internal::WireFormatLite::UInt32Size(
        this->rpcport());
  }

  // map<string, string> details = 4;
  total_size += 1 *
      ::google::protobuf::internal::FromIntSize(this->details_size());
  {
    ::google::protobuf::scoped_ptr<peer_DetailsEntry> entry;
    for (::google::protobuf::Map< ::std::string, ::std::string >::const_iterator
        it = this->details().begin();
        it != this->details().end(); ++it) {
      entry.reset(details_.NewEntryWrapper(it->first, it->second));
      total_size += ::google::protobuf::internal::WireFormatLite::
          MessageSizeNoVirtual(*entry);
    }
  }

  int cached_size = ::google::protobuf::internal::ToCachedSize(total_size);
  GOOGLE_SAFE_CONCURRENT_WRITES_BEGIN();
  _cached_size_ = cached_size;
  GOOGLE_SAFE_CONCURRENT_WRITES_END();
  return total_size;
}

void peer::MergeFrom(const ::google::protobuf::Message& from) {
// @@protoc_insertion_point(generalized_merge_from_start:pb.peer)
  if (GOOGLE_PREDICT_FALSE(&from == this)) MergeFromFail(__LINE__);
  const peer* source =
      ::google::protobuf::internal::DynamicCastToGenerated<const peer>(
          &from);
  if (source == NULL) {
  // @@protoc_insertion_point(generalized_merge_from_cast_fail:pb.peer)
    ::google::protobuf::internal::ReflectionOps::Merge(from, this);
  } else {
  // @@protoc_insertion_point(generalized_merge_from_cast_success:pb.peer)
    UnsafeMergeFrom(*source);
  }
}

void peer::MergeFrom(const peer& from) {
// @@protoc_insertion_point(class_specific_merge_from_start:pb.peer)
  if (GOOGLE_PREDICT_TRUE(&from != this)) {
    UnsafeMergeFrom(from);
  } else {
    MergeFromFail(__LINE__);
  }
}

void peer::UnsafeMergeFrom(const peer& from) {
  GOOGLE_DCHECK(&from != this);
  details_.MergeFrom(from.details_);
  if (from.ip().size() > 0) {

    ip_.AssignWithDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), from.ip_);
  }
  if (from.port() != 0) {
    set_port(from.port());
  }
  if (from.rpcport() != 0) {
    set_rpcport(from.rpcport());
  }
}

void peer::CopyFrom(const ::google::protobuf::Message& from) {
// @@protoc_insertion_point(generalized_copy_from_start:pb.peer)
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

void peer::CopyFrom(const peer& from) {
// @@protoc_insertion_point(class_specific_copy_from_start:pb.peer)
  if (&from == this) return;
  Clear();
  UnsafeMergeFrom(from);
}

bool peer::IsInitialized() const {

  return true;
}

void peer::Swap(peer* other) {
  if (other == this) return;
  InternalSwap(other);
}
void peer::InternalSwap(peer* other) {
  ip_.Swap(&other->ip_);
  std::swap(port_, other->port_);
  std::swap(rpcport_, other->rpcport_);
  details_.Swap(&other->details_);
  _internal_metadata_.Swap(&other->_internal_metadata_);
  std::swap(_cached_size_, other->_cached_size_);
}

::google::protobuf::Metadata peer::GetMetadata() const {
  protobuf_AssignDescriptorsOnce();
  ::google::protobuf::Metadata metadata;
  metadata.descriptor = peer_descriptor_;
  metadata.reflection = peer_reflection_;
  return metadata;
}

#if PROTOBUF_INLINE_NOT_IN_HEADERS
// peer

// optional string ip = 1;
void peer::clear_ip() {
  ip_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
const ::std::string& peer::ip() const {
  // @@protoc_insertion_point(field_get:pb.peer.ip)
  return ip_.GetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
void peer::set_ip(const ::std::string& value) {
  
  ip_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), value);
  // @@protoc_insertion_point(field_set:pb.peer.ip)
}
void peer::set_ip(const char* value) {
  
  ip_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::string(value));
  // @@protoc_insertion_point(field_set_char:pb.peer.ip)
}
void peer::set_ip(const char* value, size_t size) {
  
  ip_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(),
      ::std::string(reinterpret_cast<const char*>(value), size));
  // @@protoc_insertion_point(field_set_pointer:pb.peer.ip)
}
::std::string* peer::mutable_ip() {
  
  // @@protoc_insertion_point(field_mutable:pb.peer.ip)
  return ip_.MutableNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
::std::string* peer::release_ip() {
  // @@protoc_insertion_point(field_release:pb.peer.ip)
  
  return ip_.ReleaseNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
void peer::set_allocated_ip(::std::string* ip) {
  if (ip != NULL) {
    
  } else {
    
  }
  ip_.SetAllocatedNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), ip);
  // @@protoc_insertion_point(field_set_allocated:pb.peer.ip)
}

// optional uint32 port = 2;
void peer::clear_port() {
  port_ = 0u;
}
::google::protobuf::uint32 peer::port() const {
  // @@protoc_insertion_point(field_get:pb.peer.port)
  return port_;
}
void peer::set_port(::google::protobuf::uint32 value) {
  
  port_ = value;
  // @@protoc_insertion_point(field_set:pb.peer.port)
}

// optional uint32 rpcPort = 3;
void peer::clear_rpcport() {
  rpcport_ = 0u;
}
::google::protobuf::uint32 peer::rpcport() const {
  // @@protoc_insertion_point(field_get:pb.peer.rpcPort)
  return rpcport_;
}
void peer::set_rpcport(::google::protobuf::uint32 value) {
  
  rpcport_ = value;
  // @@protoc_insertion_point(field_set:pb.peer.rpcPort)
}

// map<string, string> details = 4;
int peer::details_size() const {
  return details_.size();
}
void peer::clear_details() {
  details_.Clear();
}
 const ::google::protobuf::Map< ::std::string, ::std::string >&
peer::details() const {
  // @@protoc_insertion_point(field_map:pb.peer.details)
  return details_.GetMap();
}
 ::google::protobuf::Map< ::std::string, ::std::string >*
peer::mutable_details() {
  // @@protoc_insertion_point(field_mutable_map:pb.peer.details)
  return details_.MutableMap();
}

inline const peer* peer::internal_default_instance() {
  return &peer_default_instance_.get();
}
#endif  // PROTOBUF_INLINE_NOT_IN_HEADERS

// @@protoc_insertion_point(namespace_scope)

}  // namespace pb

// @@protoc_insertion_point(global_scope)
