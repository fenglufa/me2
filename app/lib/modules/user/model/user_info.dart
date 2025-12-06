import 'package:json_annotation/json_annotation.dart';

part 'user_info.g.dart';

@JsonSerializable()
class UserInfo {
  @JsonKey(name: 'user_id')
  final int userId;
  final String nickname;
  @JsonKey(name: 'avatar_url')
  final String avatarUrl;
  @JsonKey(name: 'subscription_tier')
  final int subscriptionTier;

  UserInfo({
    required this.userId,
    required this.nickname,
    required this.avatarUrl,
    required this.subscriptionTier,
  });

  factory UserInfo.fromJson(Map<String, dynamic> json) =>
      _$UserInfoFromJson(json);
  Map<String, dynamic> toJson() => _$UserInfoToJson(this);
}
