// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: common.proto

#ifndef PROTOBUF_common_2eproto__INCLUDED
#define PROTOBUF_common_2eproto__INCLUDED

#include <string>

#include <google/protobuf/stubs/common.h>

#if GOOGLE_PROTOBUF_VERSION < 3001000
#error This file was generated by a newer version of protoc which is
#error incompatible with your Protocol Buffer headers.  Please update
#error your headers.
#endif
#if 3001000 < GOOGLE_PROTOBUF_MIN_PROTOC_VERSION
#error This file was generated by an older version of protoc which is
#error incompatible with your Protocol Buffer headers.  Please
#error regenerate this file with a newer version of protoc.
#endif

#include <google/protobuf/arena.h>
#include <google/protobuf/arenastring.h>
#include <google/protobuf/generated_message_util.h>
#include <google/protobuf/message_lite.h>
#include <google/protobuf/repeated_field.h>
#include <google/protobuf/extension_set.h>
// @@protoc_insertion_point(includes)

// Internal implementation detail -- do not call these.
void protobuf_AddDesc_common_2eproto();
void protobuf_InitDefaults_common_2eproto();
void protobuf_AssignDesc_common_2eproto();
void protobuf_ShutdownFile_common_2eproto();

class HealthReply;
class HealthRequest;

// ===================================================================

class HealthRequest : public ::google::protobuf::MessageLite /* @@protoc_insertion_point(class_definition:HealthRequest) */ {
 public:
  HealthRequest();
  virtual ~HealthRequest();

  HealthRequest(const HealthRequest& from);

  inline HealthRequest& operator=(const HealthRequest& from) {
    CopyFrom(from);
    return *this;
  }

  static const HealthRequest& default_instance();

  static const HealthRequest* internal_default_instance();

  void Swap(HealthRequest* other);

  // implements Message ----------------------------------------------

  inline HealthRequest* New() const { return New(NULL); }

  HealthRequest* New(::google::protobuf::Arena* arena) const;
  void CheckTypeAndMergeFrom(const ::google::protobuf::MessageLite& from);
  void CopyFrom(const HealthRequest& from);
  void MergeFrom(const HealthRequest& from);
  void Clear();
  bool IsInitialized() const;

  size_t ByteSizeLong() const;
  bool MergePartialFromCodedStream(
      ::google::protobuf::io::CodedInputStream* input);
  void SerializeWithCachedSizes(
      ::google::protobuf::io::CodedOutputStream* output) const;
  void DiscardUnknownFields();
  int GetCachedSize() const { return _cached_size_; }
  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const;
  void InternalSwap(HealthRequest* other);
  void UnsafeMergeFrom(const HealthRequest& from);
  private:
  inline ::google::protobuf::Arena* GetArenaNoVirtual() const {
    return _arena_ptr_;
  }
  inline ::google::protobuf::Arena* MaybeArenaPtr() const {
    return _arena_ptr_;
  }
  public:

  ::std::string GetTypeName() const;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  // @@protoc_insertion_point(class_scope:HealthRequest)
 private:

  ::google::protobuf::internal::ArenaStringPtr _unknown_fields_;
  ::google::protobuf::Arena* _arena_ptr_;

  mutable int _cached_size_;
  friend void  protobuf_InitDefaults_common_2eproto_impl();
  friend void  protobuf_AddDesc_common_2eproto_impl();
  friend void protobuf_AssignDesc_common_2eproto();
  friend void protobuf_ShutdownFile_common_2eproto();

  void InitAsDefaultInstance();
};
extern ::google::protobuf::internal::ExplicitlyConstructed<HealthRequest> HealthRequest_default_instance_;

// -------------------------------------------------------------------

class HealthReply : public ::google::protobuf::MessageLite /* @@protoc_insertion_point(class_definition:HealthReply) */ {
 public:
  HealthReply();
  virtual ~HealthReply();

  HealthReply(const HealthReply& from);

  inline HealthReply& operator=(const HealthReply& from) {
    CopyFrom(from);
    return *this;
  }

  static const HealthReply& default_instance();

  static const HealthReply* internal_default_instance();

  void Swap(HealthReply* other);

  // implements Message ----------------------------------------------

  inline HealthReply* New() const { return New(NULL); }

  HealthReply* New(::google::protobuf::Arena* arena) const;
  void CheckTypeAndMergeFrom(const ::google::protobuf::MessageLite& from);
  void CopyFrom(const HealthReply& from);
  void MergeFrom(const HealthReply& from);
  void Clear();
  bool IsInitialized() const;

  size_t ByteSizeLong() const;
  bool MergePartialFromCodedStream(
      ::google::protobuf::io::CodedInputStream* input);
  void SerializeWithCachedSizes(
      ::google::protobuf::io::CodedOutputStream* output) const;
  void DiscardUnknownFields();
  int GetCachedSize() const { return _cached_size_; }
  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const;
  void InternalSwap(HealthReply* other);
  void UnsafeMergeFrom(const HealthReply& from);
  private:
  inline ::google::protobuf::Arena* GetArenaNoVirtual() const {
    return _arena_ptr_;
  }
  inline ::google::protobuf::Arena* MaybeArenaPtr() const {
    return _arena_ptr_;
  }
  public:

  ::std::string GetTypeName() const;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  // optional bool ready = 1;
  void clear_ready();
  static const int kReadyFieldNumber = 1;
  bool ready() const;
  void set_ready(bool value);

  // optional string message = 2;
  void clear_message();
  static const int kMessageFieldNumber = 2;
  const ::std::string& message() const;
  void set_message(const ::std::string& value);
  void set_message(const char* value);
  void set_message(const char* value, size_t size);
  ::std::string* mutable_message();
  ::std::string* release_message();
  void set_allocated_message(::std::string* message);

  // @@protoc_insertion_point(class_scope:HealthReply)
 private:

  ::google::protobuf::internal::ArenaStringPtr _unknown_fields_;
  ::google::protobuf::Arena* _arena_ptr_;

  ::google::protobuf::internal::ArenaStringPtr message_;
  bool ready_;
  mutable int _cached_size_;
  friend void  protobuf_InitDefaults_common_2eproto_impl();
  friend void  protobuf_AddDesc_common_2eproto_impl();
  friend void protobuf_AssignDesc_common_2eproto();
  friend void protobuf_ShutdownFile_common_2eproto();

  void InitAsDefaultInstance();
};
extern ::google::protobuf::internal::ExplicitlyConstructed<HealthReply> HealthReply_default_instance_;

// ===================================================================


// ===================================================================

#if !PROTOBUF_INLINE_NOT_IN_HEADERS
// HealthRequest

inline const HealthRequest* HealthRequest::internal_default_instance() {
  return &HealthRequest_default_instance_.get();
}
// -------------------------------------------------------------------

// HealthReply

// optional bool ready = 1;
inline void HealthReply::clear_ready() {
  ready_ = false;
}
inline bool HealthReply::ready() const {
  // @@protoc_insertion_point(field_get:HealthReply.ready)
  return ready_;
}
inline void HealthReply::set_ready(bool value) {
  
  ready_ = value;
  // @@protoc_insertion_point(field_set:HealthReply.ready)
}

// optional string message = 2;
inline void HealthReply::clear_message() {
  message_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline const ::std::string& HealthReply::message() const {
  // @@protoc_insertion_point(field_get:HealthReply.message)
  return message_.GetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline void HealthReply::set_message(const ::std::string& value) {
  
  message_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), value);
  // @@protoc_insertion_point(field_set:HealthReply.message)
}
inline void HealthReply::set_message(const char* value) {
  
  message_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), ::std::string(value));
  // @@protoc_insertion_point(field_set_char:HealthReply.message)
}
inline void HealthReply::set_message(const char* value, size_t size) {
  
  message_.SetNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(),
      ::std::string(reinterpret_cast<const char*>(value), size));
  // @@protoc_insertion_point(field_set_pointer:HealthReply.message)
}
inline ::std::string* HealthReply::mutable_message() {
  
  // @@protoc_insertion_point(field_mutable:HealthReply.message)
  return message_.MutableNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline ::std::string* HealthReply::release_message() {
  // @@protoc_insertion_point(field_release:HealthReply.message)
  
  return message_.ReleaseNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}
inline void HealthReply::set_allocated_message(::std::string* message) {
  if (message != NULL) {
    
  } else {
    
  }
  message_.SetAllocatedNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), message);
  // @@protoc_insertion_point(field_set_allocated:HealthReply.message)
}

inline const HealthReply* HealthReply::internal_default_instance() {
  return &HealthReply_default_instance_.get();
}
#endif  // !PROTOBUF_INLINE_NOT_IN_HEADERS
// -------------------------------------------------------------------


// @@protoc_insertion_point(namespace_scope)

// @@protoc_insertion_point(global_scope)

#endif  // PROTOBUF_common_2eproto__INCLUDED
