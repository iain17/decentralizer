// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: pb/platform.proto

#include "pb/platform.pb.h"

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
// This is a temporary google only hack
#ifdef GOOGLE_PROTOBUF_ENFORCE_UNIQUENESS
#include "third_party/protobuf/version.h"
#endif
// @@protoc_insertion_point(includes)
namespace pb {
class emptyDefaultTypeInternal {
 public:
  ::google::protobuf::internal::ExplicitlyConstructed<empty>
      _instance;
} _empty_default_instance_;
class RPCHealthRequestDefaultTypeInternal {
 public:
  ::google::protobuf::internal::ExplicitlyConstructed<RPCHealthRequest>
      _instance;
} _RPCHealthRequest_default_instance_;
class RPCHealthReplyDefaultTypeInternal {
 public:
  ::google::protobuf::internal::ExplicitlyConstructed<RPCHealthReply>
      _instance;
} _RPCHealthReply_default_instance_;
}  // namespace pb
namespace protobuf_pb_2fplatform_2eproto {
void InitDefaultsemptyImpl() {
  GOOGLE_PROTOBUF_VERIFY_VERSION;

#ifdef GOOGLE_PROTOBUF_ENFORCE_UNIQUENESS
  ::google::protobuf::internal::InitProtobufDefaultsForceUnique();
#else
  ::google::protobuf::internal::InitProtobufDefaults();
#endif  // GOOGLE_PROTOBUF_ENFORCE_UNIQUENESS
  {
    void* ptr = &::pb::_empty_default_instance_;
    new (ptr) ::pb::empty();
    ::google::protobuf::internal::OnShutdownDestroyMessage(ptr);
  }
  ::pb::empty::InitAsDefaultInstance();
}

void InitDefaultsempty() {
  static GOOGLE_PROTOBUF_DECLARE_ONCE(once);
  ::google::protobuf::GoogleOnceInit(&once, &InitDefaultsemptyImpl);
}

void InitDefaultsRPCHealthRequestImpl() {
  GOOGLE_PROTOBUF_VERIFY_VERSION;

#ifdef GOOGLE_PROTOBUF_ENFORCE_UNIQUENESS
  ::google::protobuf::internal::InitProtobufDefaultsForceUnique();
#else
  ::google::protobuf::internal::InitProtobufDefaults();
#endif  // GOOGLE_PROTOBUF_ENFORCE_UNIQUENESS
  {
    void* ptr = &::pb::_RPCHealthRequest_default_instance_;
    new (ptr) ::pb::RPCHealthRequest();
    ::google::protobuf::internal::OnShutdownDestroyMessage(ptr);
  }
  ::pb::RPCHealthRequest::InitAsDefaultInstance();
}

void InitDefaultsRPCHealthRequest() {
  static GOOGLE_PROTOBUF_DECLARE_ONCE(once);
  ::google::protobuf::GoogleOnceInit(&once, &InitDefaultsRPCHealthRequestImpl);
}

void InitDefaultsRPCHealthReplyImpl() {
  GOOGLE_PROTOBUF_VERIFY_VERSION;

#ifdef GOOGLE_PROTOBUF_ENFORCE_UNIQUENESS
  ::google::protobuf::internal::InitProtobufDefaultsForceUnique();
#else
  ::google::protobuf::internal::InitProtobufDefaults();
#endif  // GOOGLE_PROTOBUF_ENFORCE_UNIQUENESS
  {
    void* ptr = &::pb::_RPCHealthReply_default_instance_;
    new (ptr) ::pb::RPCHealthReply();
    ::google::protobuf::internal::OnShutdownDestroyMessage(ptr);
  }
  ::pb::RPCHealthReply::InitAsDefaultInstance();
}

void InitDefaultsRPCHealthReply() {
  static GOOGLE_PROTOBUF_DECLARE_ONCE(once);
  ::google::protobuf::GoogleOnceInit(&once, &InitDefaultsRPCHealthReplyImpl);
}

::google::protobuf::Metadata file_level_metadata[3];

const ::google::protobuf::uint32 TableStruct::offsets[] GOOGLE_PROTOBUF_ATTRIBUTE_SECTION_VARIABLE(protodesc_cold) = {
  ~0u,  // no _has_bits_
  GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(::pb::empty, _internal_metadata_),
  ~0u,  // no _extensions_
  ~0u,  // no _oneof_case_
  ~0u,  // no _weak_field_map_
  ~0u,  // no _has_bits_
  GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(::pb::RPCHealthRequest, _internal_metadata_),
  ~0u,  // no _extensions_
  ~0u,  // no _oneof_case_
  ~0u,  // no _weak_field_map_
  ~0u,  // no _has_bits_
  GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(::pb::RPCHealthReply, _internal_metadata_),
  ~0u,  // no _extensions_
  ~0u,  // no _oneof_case_
  ~0u,  // no _weak_field_map_
  GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(::pb::RPCHealthReply, ready_),
  GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(::pb::RPCHealthReply, message_),
  GOOGLE_PROTOBUF_GENERATED_MESSAGE_FIELD_OFFSET(::pb::RPCHealthReply, basepath_),
};
static const ::google::protobuf::internal::MigrationSchema schemas[] GOOGLE_PROTOBUF_ATTRIBUTE_SECTION_VARIABLE(protodesc_cold) = {
  { 0, -1, sizeof(::pb::empty)},
  { 5, -1, sizeof(::pb::RPCHealthRequest)},
  { 10, -1, sizeof(::pb::RPCHealthReply)},
};

static ::google::protobuf::Message const * const file_default_instances[] = {
  reinterpret_cast<const ::google::protobuf::Message*>(&::pb::_empty_default_instance_),
  reinterpret_cast<const ::google::protobuf::Message*>(&::pb::_RPCHealthRequest_default_instance_),
  reinterpret_cast<const ::google::protobuf::Message*>(&::pb::_RPCHealthReply_default_instance_),
};

void protobuf_AssignDescriptors() {
  AddDescriptors();
  ::google::protobuf::MessageFactory* factory = NULL;
  AssignDescriptors(
      "pb/platform.proto", schemas, file_default_instances, TableStruct::offsets, factory,
      file_level_metadata, NULL, NULL);
}

void protobuf_AssignDescriptorsOnce() {
  static GOOGLE_PROTOBUF_DECLARE_ONCE(once);
  ::google::protobuf::GoogleOnceInit(&once, &protobuf_AssignDescriptors);
}

void protobuf_RegisterTypes(const ::std::string&) GOOGLE_PROTOBUF_ATTRIBUTE_COLD;
void protobuf_RegisterTypes(const ::std::string&) {
  protobuf_AssignDescriptorsOnce();
  ::google::protobuf::internal::RegisterAllTypes(file_level_metadata, 3);
}

void AddDescriptorsImpl() {
  InitDefaults();
  static const char descriptor[] GOOGLE_PROTOBUF_ATTRIBUTE_SECTION_VARIABLE(protodesc_cold) = {
      "\n\021pb/platform.proto\022\002pb\032\024pb/matchmaking."
      "proto\032\024pb/addressbook.proto\032\020pb/storage."
      "proto\032\022pb/messaging.proto\032\022pb/publisher."
      "proto\032\034google/api/annotations.proto\"\007\n\005e"
      "mpty\"\022\n\020RPCHealthRequest\"B\n\016RPCHealthRep"
      "ly\022\r\n\005ready\030\001 \001(\010\022\017\n\007message\030\002 \001(\t\022\020\n\010ba"
      "sePath\030\003 \001(\t2\243\013\n\rDecentralizer\022I\n\tGetHea"
      "lth\022\024.pb.RPCHealthRequest\032\022.pb.RPCHealth"
      "Reply\"\022\202\323\344\223\002\014\022\n/v1/health\022_\n\rUpsertSessi"
      "on\022\033.pb.RPCUpsertSessionRequest\032\034.pb.RPC"
      "UpsertSessionResponse\"\023\202\323\344\223\002\r\"\013/v1/sessi"
      "on\022L\n\rDeleteSession\022\033.pb.RPCDeleteSessio"
      "nRequest\032\034.pb.RPCDeleteSessionResponse\"\000"
      "\022\207\001\n\026GetSessionIdsByDetails\022$.pb.RPCGetS"
      "essionIdsByDetailsRequest\032\034.pb.RPCGetSes"
      "sionIdsResponse\")\202\323\344\223\002#\022!/v1/sessions/{t"
      "ype}/{key}/{value}\022^\n\026GetSessionIdsByPee"
      "rIds\022$.pb.RPCGetSessionIdsByPeerIdsReque"
      "st\032\034.pb.RPCGetSessionIdsResponse\"\000\022b\n\nGe"
      "tSession\022\030.pb.RPCGetSessionRequest\032\031.pb."
      "RPCGetSessionResponse\"\037\202\323\344\223\002\031\022\027/v1/sessi"
      "on/{sessionId}\022C\n\nUpsertPeer\022\030.pb.RPCUps"
      "ertPeerRequest\032\031.pb.RPCUpsertPeerRespons"
      "e\"\000\022b\n\nGetPeerIds\022\030.pb.RPCGetPeerIdsRequ"
      "est\032\031.pb.RPCGetPeerIdsResponse\"\037\202\323\344\223\002\031\022\027"
      "/v1/peers/{key}/{value}\022P\n\007GetPeer\022\025.pb."
      "RPCGetPeerRequest\032\026.pb.RPCGetPeerRespons"
      "e\"\026\202\323\344\223\002\020\022\016/v1/peer/{pId}\022L\n\rWritePeerFi"
      "le\022\033.pb.RPCWritePeerFileRequest\032\034.pb.RPC"
      "WritePeerFileResponse\"\000\022F\n\013GetPeerFile\022\031"
      ".pb.RPCGetPeerFileRequest\032\032.pb.RPCGetPee"
      "rFileResponse\"\000\0226\n\021SendDirectMessage\022\024.p"
      "b.RPCDirectMessage\032\t.pb.empty\"\000\022T\n\024Recei"
      "veDirectMessage\022\".pb.RPCReceiveDirectMes"
      "sageRequest\032\024.pb.RPCDirectMessage\"\0000\001\022K\n"
      "\026publishPublisherUpdate\022$.pb.RPCPublishP"
      "ublisherUpdateRequest\032\t.pb.empty\"\000\022k\n\026Ge"
      "tPublisherDefinition\022!.pb.GetPublisherDe"
      "finitionRequest\032\027.pb.PublisherDefinition"
      "\"\025\202\323\344\223\002\017\022\r/v1/publisher\022q\n\020GetPublisherF"
      "ile\022\036.pb.RPCGetPublisherFileRequest\032\037.pb"
      ".RPCGetPublisherFileResponse\"\034\202\323\344\223\002\026\022\024/v"
      "1/publisher/{name}b\006proto3"
  };
  ::google::protobuf::DescriptorPool::InternalAddGeneratedFile(
      descriptor, 1706);
  ::google::protobuf::MessageFactory::InternalRegisterGeneratedFile(
    "pb/platform.proto", &protobuf_RegisterTypes);
  ::protobuf_pb_2fmatchmaking_2eproto::AddDescriptors();
  ::protobuf_pb_2faddressbook_2eproto::AddDescriptors();
  ::protobuf_pb_2fstorage_2eproto::AddDescriptors();
  ::protobuf_pb_2fmessaging_2eproto::AddDescriptors();
  ::protobuf_pb_2fpublisher_2eproto::AddDescriptors();
}

void AddDescriptors() {
  static GOOGLE_PROTOBUF_DECLARE_ONCE(once);
  ::google::protobuf::GoogleOnceInit(&once, &AddDescriptorsImpl);
}
// Force AddDescriptors() to be called at dynamic initialization time.
struct StaticDescriptorInitializer {
  StaticDescriptorInitializer() {
    AddDescriptors();
  }
} static_descriptor_initializer;
}  // namespace protobuf_pb_2fplatform_2eproto
namespace pb {

// ===================================================================

void empty::InitAsDefaultInstance() {
}
#if !defined(_MSC_VER) || _MSC_VER >= 1900
#endif  // !defined(_MSC_VER) || _MSC_VER >= 1900

empty::empty()
  : ::google::protobuf::Message(), _internal_metadata_(NULL) {
  if (GOOGLE_PREDICT_TRUE(this != internal_default_instance())) {
    ::protobuf_pb_2fplatform_2eproto::InitDefaultsempty();
  }
  SharedCtor();
  // @@protoc_insertion_point(constructor:pb.empty)
}
empty::empty(const empty& from)
  : ::google::protobuf::Message(),
      _internal_metadata_(NULL),
      _cached_size_(0) {
  _internal_metadata_.MergeFrom(from._internal_metadata_);
  // @@protoc_insertion_point(copy_constructor:pb.empty)
}

void empty::SharedCtor() {
  _cached_size_ = 0;
}

empty::~empty() {
  // @@protoc_insertion_point(destructor:pb.empty)
  SharedDtor();
}

void empty::SharedDtor() {
}

void empty::SetCachedSize(int size) const {
  GOOGLE_SAFE_CONCURRENT_WRITES_BEGIN();
  _cached_size_ = size;
  GOOGLE_SAFE_CONCURRENT_WRITES_END();
}
const ::google::protobuf::Descriptor* empty::descriptor() {
  ::protobuf_pb_2fplatform_2eproto::protobuf_AssignDescriptorsOnce();
  return ::protobuf_pb_2fplatform_2eproto::file_level_metadata[kIndexInFileMessages].descriptor;
}

const empty& empty::default_instance() {
  ::protobuf_pb_2fplatform_2eproto::InitDefaultsempty();
  return *internal_default_instance();
}

empty* empty::New(::google::protobuf::Arena* arena) const {
  empty* n = new empty;
  if (arena != NULL) {
    arena->Own(n);
  }
  return n;
}

void empty::Clear() {
// @@protoc_insertion_point(message_clear_start:pb.empty)
  ::google::protobuf::uint32 cached_has_bits = 0;
  // Prevent compiler warnings about cached_has_bits being unused
  (void) cached_has_bits;

  _internal_metadata_.Clear();
}

bool empty::MergePartialFromCodedStream(
    ::google::protobuf::io::CodedInputStream* input) {
#define DO_(EXPRESSION) if (!GOOGLE_PREDICT_TRUE(EXPRESSION)) goto failure
  ::google::protobuf::uint32 tag;
  // @@protoc_insertion_point(parse_start:pb.empty)
  for (;;) {
    ::std::pair< ::google::protobuf::uint32, bool> p = input->ReadTagWithCutoffNoLastTag(127u);
    tag = p.first;
    if (!p.second) goto handle_unusual;
  handle_unusual:
    if (tag == 0) {
      goto success;
    }
    DO_(::google::protobuf::internal::WireFormat::SkipField(
          input, tag, _internal_metadata_.mutable_unknown_fields()));
  }
success:
  // @@protoc_insertion_point(parse_success:pb.empty)
  return true;
failure:
  // @@protoc_insertion_point(parse_failure:pb.empty)
  return false;
#undef DO_
}

void empty::SerializeWithCachedSizes(
    ::google::protobuf::io::CodedOutputStream* output) const {
  // @@protoc_insertion_point(serialize_start:pb.empty)
  ::google::protobuf::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  if ((_internal_metadata_.have_unknown_fields() &&  ::google::protobuf::internal::GetProto3PreserveUnknownsDefault())) {
    ::google::protobuf::internal::WireFormat::SerializeUnknownFields(
        (::google::protobuf::internal::GetProto3PreserveUnknownsDefault()   ? _internal_metadata_.unknown_fields()   : _internal_metadata_.default_instance()), output);
  }
  // @@protoc_insertion_point(serialize_end:pb.empty)
}

::google::protobuf::uint8* empty::InternalSerializeWithCachedSizesToArray(
    bool deterministic, ::google::protobuf::uint8* target) const {
  (void)deterministic; // Unused
  // @@protoc_insertion_point(serialize_to_array_start:pb.empty)
  ::google::protobuf::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  if ((_internal_metadata_.have_unknown_fields() &&  ::google::protobuf::internal::GetProto3PreserveUnknownsDefault())) {
    target = ::google::protobuf::internal::WireFormat::SerializeUnknownFieldsToArray(
        (::google::protobuf::internal::GetProto3PreserveUnknownsDefault()   ? _internal_metadata_.unknown_fields()   : _internal_metadata_.default_instance()), target);
  }
  // @@protoc_insertion_point(serialize_to_array_end:pb.empty)
  return target;
}

size_t empty::ByteSizeLong() const {
// @@protoc_insertion_point(message_byte_size_start:pb.empty)
  size_t total_size = 0;

  if ((_internal_metadata_.have_unknown_fields() &&  ::google::protobuf::internal::GetProto3PreserveUnknownsDefault())) {
    total_size +=
      ::google::protobuf::internal::WireFormat::ComputeUnknownFieldsSize(
        (::google::protobuf::internal::GetProto3PreserveUnknownsDefault()   ? _internal_metadata_.unknown_fields()   : _internal_metadata_.default_instance()));
  }
  int cached_size = ::google::protobuf::internal::ToCachedSize(total_size);
  GOOGLE_SAFE_CONCURRENT_WRITES_BEGIN();
  _cached_size_ = cached_size;
  GOOGLE_SAFE_CONCURRENT_WRITES_END();
  return total_size;
}

void empty::MergeFrom(const ::google::protobuf::Message& from) {
// @@protoc_insertion_point(generalized_merge_from_start:pb.empty)
  GOOGLE_DCHECK_NE(&from, this);
  const empty* source =
      ::google::protobuf::internal::DynamicCastToGenerated<const empty>(
          &from);
  if (source == NULL) {
  // @@protoc_insertion_point(generalized_merge_from_cast_fail:pb.empty)
    ::google::protobuf::internal::ReflectionOps::Merge(from, this);
  } else {
  // @@protoc_insertion_point(generalized_merge_from_cast_success:pb.empty)
    MergeFrom(*source);
  }
}

void empty::MergeFrom(const empty& from) {
// @@protoc_insertion_point(class_specific_merge_from_start:pb.empty)
  GOOGLE_DCHECK_NE(&from, this);
  _internal_metadata_.MergeFrom(from._internal_metadata_);
  ::google::protobuf::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

}

void empty::CopyFrom(const ::google::protobuf::Message& from) {
// @@protoc_insertion_point(generalized_copy_from_start:pb.empty)
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

void empty::CopyFrom(const empty& from) {
// @@protoc_insertion_point(class_specific_copy_from_start:pb.empty)
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

bool empty::IsInitialized() const {
  return true;
}

void empty::Swap(empty* other) {
  if (other == this) return;
  InternalSwap(other);
}
void empty::InternalSwap(empty* other) {
  using std::swap;
  _internal_metadata_.Swap(&other->_internal_metadata_);
  swap(_cached_size_, other->_cached_size_);
}

::google::protobuf::Metadata empty::GetMetadata() const {
  protobuf_pb_2fplatform_2eproto::protobuf_AssignDescriptorsOnce();
  return ::protobuf_pb_2fplatform_2eproto::file_level_metadata[kIndexInFileMessages];
}


// ===================================================================

void RPCHealthRequest::InitAsDefaultInstance() {
}
#if !defined(_MSC_VER) || _MSC_VER >= 1900
#endif  // !defined(_MSC_VER) || _MSC_VER >= 1900

RPCHealthRequest::RPCHealthRequest()
  : ::google::protobuf::Message(), _internal_metadata_(NULL) {
  if (GOOGLE_PREDICT_TRUE(this != internal_default_instance())) {
    ::protobuf_pb_2fplatform_2eproto::InitDefaultsRPCHealthRequest();
  }
  SharedCtor();
  // @@protoc_insertion_point(constructor:pb.RPCHealthRequest)
}
RPCHealthRequest::RPCHealthRequest(const RPCHealthRequest& from)
  : ::google::protobuf::Message(),
      _internal_metadata_(NULL),
      _cached_size_(0) {
  _internal_metadata_.MergeFrom(from._internal_metadata_);
  // @@protoc_insertion_point(copy_constructor:pb.RPCHealthRequest)
}

void RPCHealthRequest::SharedCtor() {
  _cached_size_ = 0;
}

RPCHealthRequest::~RPCHealthRequest() {
  // @@protoc_insertion_point(destructor:pb.RPCHealthRequest)
  SharedDtor();
}

void RPCHealthRequest::SharedDtor() {
}

void RPCHealthRequest::SetCachedSize(int size) const {
  GOOGLE_SAFE_CONCURRENT_WRITES_BEGIN();
  _cached_size_ = size;
  GOOGLE_SAFE_CONCURRENT_WRITES_END();
}
const ::google::protobuf::Descriptor* RPCHealthRequest::descriptor() {
  ::protobuf_pb_2fplatform_2eproto::protobuf_AssignDescriptorsOnce();
  return ::protobuf_pb_2fplatform_2eproto::file_level_metadata[kIndexInFileMessages].descriptor;
}

const RPCHealthRequest& RPCHealthRequest::default_instance() {
  ::protobuf_pb_2fplatform_2eproto::InitDefaultsRPCHealthRequest();
  return *internal_default_instance();
}

RPCHealthRequest* RPCHealthRequest::New(::google::protobuf::Arena* arena) const {
  RPCHealthRequest* n = new RPCHealthRequest;
  if (arena != NULL) {
    arena->Own(n);
  }
  return n;
}

void RPCHealthRequest::Clear() {
// @@protoc_insertion_point(message_clear_start:pb.RPCHealthRequest)
  ::google::protobuf::uint32 cached_has_bits = 0;
  // Prevent compiler warnings about cached_has_bits being unused
  (void) cached_has_bits;

  _internal_metadata_.Clear();
}

bool RPCHealthRequest::MergePartialFromCodedStream(
    ::google::protobuf::io::CodedInputStream* input) {
#define DO_(EXPRESSION) if (!GOOGLE_PREDICT_TRUE(EXPRESSION)) goto failure
  ::google::protobuf::uint32 tag;
  // @@protoc_insertion_point(parse_start:pb.RPCHealthRequest)
  for (;;) {
    ::std::pair< ::google::protobuf::uint32, bool> p = input->ReadTagWithCutoffNoLastTag(127u);
    tag = p.first;
    if (!p.second) goto handle_unusual;
  handle_unusual:
    if (tag == 0) {
      goto success;
    }
    DO_(::google::protobuf::internal::WireFormat::SkipField(
          input, tag, _internal_metadata_.mutable_unknown_fields()));
  }
success:
  // @@protoc_insertion_point(parse_success:pb.RPCHealthRequest)
  return true;
failure:
  // @@protoc_insertion_point(parse_failure:pb.RPCHealthRequest)
  return false;
#undef DO_
}

void RPCHealthRequest::SerializeWithCachedSizes(
    ::google::protobuf::io::CodedOutputStream* output) const {
  // @@protoc_insertion_point(serialize_start:pb.RPCHealthRequest)
  ::google::protobuf::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  if ((_internal_metadata_.have_unknown_fields() &&  ::google::protobuf::internal::GetProto3PreserveUnknownsDefault())) {
    ::google::protobuf::internal::WireFormat::SerializeUnknownFields(
        (::google::protobuf::internal::GetProto3PreserveUnknownsDefault()   ? _internal_metadata_.unknown_fields()   : _internal_metadata_.default_instance()), output);
  }
  // @@protoc_insertion_point(serialize_end:pb.RPCHealthRequest)
}

::google::protobuf::uint8* RPCHealthRequest::InternalSerializeWithCachedSizesToArray(
    bool deterministic, ::google::protobuf::uint8* target) const {
  (void)deterministic; // Unused
  // @@protoc_insertion_point(serialize_to_array_start:pb.RPCHealthRequest)
  ::google::protobuf::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  if ((_internal_metadata_.have_unknown_fields() &&  ::google::protobuf::internal::GetProto3PreserveUnknownsDefault())) {
    target = ::google::protobuf::internal::WireFormat::SerializeUnknownFieldsToArray(
        (::google::protobuf::internal::GetProto3PreserveUnknownsDefault()   ? _internal_metadata_.unknown_fields()   : _internal_metadata_.default_instance()), target);
  }
  // @@protoc_insertion_point(serialize_to_array_end:pb.RPCHealthRequest)
  return target;
}

size_t RPCHealthRequest::ByteSizeLong() const {
// @@protoc_insertion_point(message_byte_size_start:pb.RPCHealthRequest)
  size_t total_size = 0;

  if ((_internal_metadata_.have_unknown_fields() &&  ::google::protobuf::internal::GetProto3PreserveUnknownsDefault())) {
    total_size +=
      ::google::protobuf::internal::WireFormat::ComputeUnknownFieldsSize(
        (::google::protobuf::internal::GetProto3PreserveUnknownsDefault()   ? _internal_metadata_.unknown_fields()   : _internal_metadata_.default_instance()));
  }
  int cached_size = ::google::protobuf::internal::ToCachedSize(total_size);
  GOOGLE_SAFE_CONCURRENT_WRITES_BEGIN();
  _cached_size_ = cached_size;
  GOOGLE_SAFE_CONCURRENT_WRITES_END();
  return total_size;
}

void RPCHealthRequest::MergeFrom(const ::google::protobuf::Message& from) {
// @@protoc_insertion_point(generalized_merge_from_start:pb.RPCHealthRequest)
  GOOGLE_DCHECK_NE(&from, this);
  const RPCHealthRequest* source =
      ::google::protobuf::internal::DynamicCastToGenerated<const RPCHealthRequest>(
          &from);
  if (source == NULL) {
  // @@protoc_insertion_point(generalized_merge_from_cast_fail:pb.RPCHealthRequest)
    ::google::protobuf::internal::ReflectionOps::Merge(from, this);
  } else {
  // @@protoc_insertion_point(generalized_merge_from_cast_success:pb.RPCHealthRequest)
    MergeFrom(*source);
  }
}

void RPCHealthRequest::MergeFrom(const RPCHealthRequest& from) {
// @@protoc_insertion_point(class_specific_merge_from_start:pb.RPCHealthRequest)
  GOOGLE_DCHECK_NE(&from, this);
  _internal_metadata_.MergeFrom(from._internal_metadata_);
  ::google::protobuf::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

}

void RPCHealthRequest::CopyFrom(const ::google::protobuf::Message& from) {
// @@protoc_insertion_point(generalized_copy_from_start:pb.RPCHealthRequest)
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

void RPCHealthRequest::CopyFrom(const RPCHealthRequest& from) {
// @@protoc_insertion_point(class_specific_copy_from_start:pb.RPCHealthRequest)
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

bool RPCHealthRequest::IsInitialized() const {
  return true;
}

void RPCHealthRequest::Swap(RPCHealthRequest* other) {
  if (other == this) return;
  InternalSwap(other);
}
void RPCHealthRequest::InternalSwap(RPCHealthRequest* other) {
  using std::swap;
  _internal_metadata_.Swap(&other->_internal_metadata_);
  swap(_cached_size_, other->_cached_size_);
}

::google::protobuf::Metadata RPCHealthRequest::GetMetadata() const {
  protobuf_pb_2fplatform_2eproto::protobuf_AssignDescriptorsOnce();
  return ::protobuf_pb_2fplatform_2eproto::file_level_metadata[kIndexInFileMessages];
}


// ===================================================================

void RPCHealthReply::InitAsDefaultInstance() {
}
#if !defined(_MSC_VER) || _MSC_VER >= 1900
const int RPCHealthReply::kReadyFieldNumber;
const int RPCHealthReply::kMessageFieldNumber;
const int RPCHealthReply::kBasePathFieldNumber;
#endif  // !defined(_MSC_VER) || _MSC_VER >= 1900

RPCHealthReply::RPCHealthReply()
  : ::google::protobuf::Message(), _internal_metadata_(NULL) {
  if (GOOGLE_PREDICT_TRUE(this != internal_default_instance())) {
    ::protobuf_pb_2fplatform_2eproto::InitDefaultsRPCHealthReply();
  }
  SharedCtor();
  // @@protoc_insertion_point(constructor:pb.RPCHealthReply)
}
RPCHealthReply::RPCHealthReply(const RPCHealthReply& from)
  : ::google::protobuf::Message(),
      _internal_metadata_(NULL),
      _cached_size_(0) {
  _internal_metadata_.MergeFrom(from._internal_metadata_);
  message_.UnsafeSetDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
  if (from.message().size() > 0) {
    message_.AssignWithDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), from.message_);
  }
  basepath_.UnsafeSetDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
  if (from.basepath().size() > 0) {
    basepath_.AssignWithDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), from.basepath_);
  }
  ready_ = from.ready_;
  // @@protoc_insertion_point(copy_constructor:pb.RPCHealthReply)
}

void RPCHealthReply::SharedCtor() {
  message_.UnsafeSetDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
  basepath_.UnsafeSetDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
  ready_ = false;
  _cached_size_ = 0;
}

RPCHealthReply::~RPCHealthReply() {
  // @@protoc_insertion_point(destructor:pb.RPCHealthReply)
  SharedDtor();
}

void RPCHealthReply::SharedDtor() {
  message_.DestroyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
  basepath_.DestroyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
}

void RPCHealthReply::SetCachedSize(int size) const {
  GOOGLE_SAFE_CONCURRENT_WRITES_BEGIN();
  _cached_size_ = size;
  GOOGLE_SAFE_CONCURRENT_WRITES_END();
}
const ::google::protobuf::Descriptor* RPCHealthReply::descriptor() {
  ::protobuf_pb_2fplatform_2eproto::protobuf_AssignDescriptorsOnce();
  return ::protobuf_pb_2fplatform_2eproto::file_level_metadata[kIndexInFileMessages].descriptor;
}

const RPCHealthReply& RPCHealthReply::default_instance() {
  ::protobuf_pb_2fplatform_2eproto::InitDefaultsRPCHealthReply();
  return *internal_default_instance();
}

RPCHealthReply* RPCHealthReply::New(::google::protobuf::Arena* arena) const {
  RPCHealthReply* n = new RPCHealthReply;
  if (arena != NULL) {
    arena->Own(n);
  }
  return n;
}

void RPCHealthReply::Clear() {
// @@protoc_insertion_point(message_clear_start:pb.RPCHealthReply)
  ::google::protobuf::uint32 cached_has_bits = 0;
  // Prevent compiler warnings about cached_has_bits being unused
  (void) cached_has_bits;

  message_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
  basepath_.ClearToEmptyNoArena(&::google::protobuf::internal::GetEmptyStringAlreadyInited());
  ready_ = false;
  _internal_metadata_.Clear();
}

bool RPCHealthReply::MergePartialFromCodedStream(
    ::google::protobuf::io::CodedInputStream* input) {
#define DO_(EXPRESSION) if (!GOOGLE_PREDICT_TRUE(EXPRESSION)) goto failure
  ::google::protobuf::uint32 tag;
  // @@protoc_insertion_point(parse_start:pb.RPCHealthReply)
  for (;;) {
    ::std::pair< ::google::protobuf::uint32, bool> p = input->ReadTagWithCutoffNoLastTag(127u);
    tag = p.first;
    if (!p.second) goto handle_unusual;
    switch (::google::protobuf::internal::WireFormatLite::GetTagFieldNumber(tag)) {
      // bool ready = 1;
      case 1: {
        if (static_cast< ::google::protobuf::uint8>(tag) ==
            static_cast< ::google::protobuf::uint8>(8u /* 8 & 0xFF */)) {

          DO_((::google::protobuf::internal::WireFormatLite::ReadPrimitive<
                   bool, ::google::protobuf::internal::WireFormatLite::TYPE_BOOL>(
                 input, &ready_)));
        } else {
          goto handle_unusual;
        }
        break;
      }

      // string message = 2;
      case 2: {
        if (static_cast< ::google::protobuf::uint8>(tag) ==
            static_cast< ::google::protobuf::uint8>(18u /* 18 & 0xFF */)) {
          DO_(::google::protobuf::internal::WireFormatLite::ReadString(
                input, this->mutable_message()));
          DO_(::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
            this->message().data(), static_cast<int>(this->message().length()),
            ::google::protobuf::internal::WireFormatLite::PARSE,
            "pb.RPCHealthReply.message"));
        } else {
          goto handle_unusual;
        }
        break;
      }

      // string basePath = 3;
      case 3: {
        if (static_cast< ::google::protobuf::uint8>(tag) ==
            static_cast< ::google::protobuf::uint8>(26u /* 26 & 0xFF */)) {
          DO_(::google::protobuf::internal::WireFormatLite::ReadString(
                input, this->mutable_basepath()));
          DO_(::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
            this->basepath().data(), static_cast<int>(this->basepath().length()),
            ::google::protobuf::internal::WireFormatLite::PARSE,
            "pb.RPCHealthReply.basePath"));
        } else {
          goto handle_unusual;
        }
        break;
      }

      default: {
      handle_unusual:
        if (tag == 0) {
          goto success;
        }
        DO_(::google::protobuf::internal::WireFormat::SkipField(
              input, tag, _internal_metadata_.mutable_unknown_fields()));
        break;
      }
    }
  }
success:
  // @@protoc_insertion_point(parse_success:pb.RPCHealthReply)
  return true;
failure:
  // @@protoc_insertion_point(parse_failure:pb.RPCHealthReply)
  return false;
#undef DO_
}

void RPCHealthReply::SerializeWithCachedSizes(
    ::google::protobuf::io::CodedOutputStream* output) const {
  // @@protoc_insertion_point(serialize_start:pb.RPCHealthReply)
  ::google::protobuf::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  // bool ready = 1;
  if (this->ready() != 0) {
    ::google::protobuf::internal::WireFormatLite::WriteBool(1, this->ready(), output);
  }

  // string message = 2;
  if (this->message().size() > 0) {
    ::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
      this->message().data(), static_cast<int>(this->message().length()),
      ::google::protobuf::internal::WireFormatLite::SERIALIZE,
      "pb.RPCHealthReply.message");
    ::google::protobuf::internal::WireFormatLite::WriteStringMaybeAliased(
      2, this->message(), output);
  }

  // string basePath = 3;
  if (this->basepath().size() > 0) {
    ::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
      this->basepath().data(), static_cast<int>(this->basepath().length()),
      ::google::protobuf::internal::WireFormatLite::SERIALIZE,
      "pb.RPCHealthReply.basePath");
    ::google::protobuf::internal::WireFormatLite::WriteStringMaybeAliased(
      3, this->basepath(), output);
  }

  if ((_internal_metadata_.have_unknown_fields() &&  ::google::protobuf::internal::GetProto3PreserveUnknownsDefault())) {
    ::google::protobuf::internal::WireFormat::SerializeUnknownFields(
        (::google::protobuf::internal::GetProto3PreserveUnknownsDefault()   ? _internal_metadata_.unknown_fields()   : _internal_metadata_.default_instance()), output);
  }
  // @@protoc_insertion_point(serialize_end:pb.RPCHealthReply)
}

::google::protobuf::uint8* RPCHealthReply::InternalSerializeWithCachedSizesToArray(
    bool deterministic, ::google::protobuf::uint8* target) const {
  (void)deterministic; // Unused
  // @@protoc_insertion_point(serialize_to_array_start:pb.RPCHealthReply)
  ::google::protobuf::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  // bool ready = 1;
  if (this->ready() != 0) {
    target = ::google::protobuf::internal::WireFormatLite::WriteBoolToArray(1, this->ready(), target);
  }

  // string message = 2;
  if (this->message().size() > 0) {
    ::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
      this->message().data(), static_cast<int>(this->message().length()),
      ::google::protobuf::internal::WireFormatLite::SERIALIZE,
      "pb.RPCHealthReply.message");
    target =
      ::google::protobuf::internal::WireFormatLite::WriteStringToArray(
        2, this->message(), target);
  }

  // string basePath = 3;
  if (this->basepath().size() > 0) {
    ::google::protobuf::internal::WireFormatLite::VerifyUtf8String(
      this->basepath().data(), static_cast<int>(this->basepath().length()),
      ::google::protobuf::internal::WireFormatLite::SERIALIZE,
      "pb.RPCHealthReply.basePath");
    target =
      ::google::protobuf::internal::WireFormatLite::WriteStringToArray(
        3, this->basepath(), target);
  }

  if ((_internal_metadata_.have_unknown_fields() &&  ::google::protobuf::internal::GetProto3PreserveUnknownsDefault())) {
    target = ::google::protobuf::internal::WireFormat::SerializeUnknownFieldsToArray(
        (::google::protobuf::internal::GetProto3PreserveUnknownsDefault()   ? _internal_metadata_.unknown_fields()   : _internal_metadata_.default_instance()), target);
  }
  // @@protoc_insertion_point(serialize_to_array_end:pb.RPCHealthReply)
  return target;
}

size_t RPCHealthReply::ByteSizeLong() const {
// @@protoc_insertion_point(message_byte_size_start:pb.RPCHealthReply)
  size_t total_size = 0;

  if ((_internal_metadata_.have_unknown_fields() &&  ::google::protobuf::internal::GetProto3PreserveUnknownsDefault())) {
    total_size +=
      ::google::protobuf::internal::WireFormat::ComputeUnknownFieldsSize(
        (::google::protobuf::internal::GetProto3PreserveUnknownsDefault()   ? _internal_metadata_.unknown_fields()   : _internal_metadata_.default_instance()));
  }
  // string message = 2;
  if (this->message().size() > 0) {
    total_size += 1 +
      ::google::protobuf::internal::WireFormatLite::StringSize(
        this->message());
  }

  // string basePath = 3;
  if (this->basepath().size() > 0) {
    total_size += 1 +
      ::google::protobuf::internal::WireFormatLite::StringSize(
        this->basepath());
  }

  // bool ready = 1;
  if (this->ready() != 0) {
    total_size += 1 + 1;
  }

  int cached_size = ::google::protobuf::internal::ToCachedSize(total_size);
  GOOGLE_SAFE_CONCURRENT_WRITES_BEGIN();
  _cached_size_ = cached_size;
  GOOGLE_SAFE_CONCURRENT_WRITES_END();
  return total_size;
}

void RPCHealthReply::MergeFrom(const ::google::protobuf::Message& from) {
// @@protoc_insertion_point(generalized_merge_from_start:pb.RPCHealthReply)
  GOOGLE_DCHECK_NE(&from, this);
  const RPCHealthReply* source =
      ::google::protobuf::internal::DynamicCastToGenerated<const RPCHealthReply>(
          &from);
  if (source == NULL) {
  // @@protoc_insertion_point(generalized_merge_from_cast_fail:pb.RPCHealthReply)
    ::google::protobuf::internal::ReflectionOps::Merge(from, this);
  } else {
  // @@protoc_insertion_point(generalized_merge_from_cast_success:pb.RPCHealthReply)
    MergeFrom(*source);
  }
}

void RPCHealthReply::MergeFrom(const RPCHealthReply& from) {
// @@protoc_insertion_point(class_specific_merge_from_start:pb.RPCHealthReply)
  GOOGLE_DCHECK_NE(&from, this);
  _internal_metadata_.MergeFrom(from._internal_metadata_);
  ::google::protobuf::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  if (from.message().size() > 0) {

    message_.AssignWithDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), from.message_);
  }
  if (from.basepath().size() > 0) {

    basepath_.AssignWithDefault(&::google::protobuf::internal::GetEmptyStringAlreadyInited(), from.basepath_);
  }
  if (from.ready() != 0) {
    set_ready(from.ready());
  }
}

void RPCHealthReply::CopyFrom(const ::google::protobuf::Message& from) {
// @@protoc_insertion_point(generalized_copy_from_start:pb.RPCHealthReply)
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

void RPCHealthReply::CopyFrom(const RPCHealthReply& from) {
// @@protoc_insertion_point(class_specific_copy_from_start:pb.RPCHealthReply)
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

bool RPCHealthReply::IsInitialized() const {
  return true;
}

void RPCHealthReply::Swap(RPCHealthReply* other) {
  if (other == this) return;
  InternalSwap(other);
}
void RPCHealthReply::InternalSwap(RPCHealthReply* other) {
  using std::swap;
  message_.Swap(&other->message_);
  basepath_.Swap(&other->basepath_);
  swap(ready_, other->ready_);
  _internal_metadata_.Swap(&other->_internal_metadata_);
  swap(_cached_size_, other->_cached_size_);
}

::google::protobuf::Metadata RPCHealthReply::GetMetadata() const {
  protobuf_pb_2fplatform_2eproto::protobuf_AssignDescriptorsOnce();
  return ::protobuf_pb_2fplatform_2eproto::file_level_metadata[kIndexInFileMessages];
}


// @@protoc_insertion_point(namespace_scope)
}  // namespace pb

// @@protoc_insertion_point(global_scope)
