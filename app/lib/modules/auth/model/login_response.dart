import 'package:json_annotation/json_annotation.dart';

part 'login_response.g.dart';

@JsonSerializable()
class LoginResponse {
  final String token;
  @JsonKey(name: 'user_id')
  final int userId;
  @JsonKey(name: 'avatar_id')
  final int? avatarId;

  LoginResponse({
    required this.token,
    required this.userId,
    this.avatarId,
  });

  factory LoginResponse.fromJson(Map<String, dynamic> json) =>
      _$LoginResponseFromJson(json);

  Map<String, dynamic> toJson() => _$LoginResponseToJson(this);
}
