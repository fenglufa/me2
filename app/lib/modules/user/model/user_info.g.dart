// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'user_info.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

UserInfo _$UserInfoFromJson(Map<String, dynamic> json) => UserInfo(
      userId: (json['user_id'] as num).toInt(),
      nickname: json['nickname'] as String,
      avatarUrl: json['avatar_url'] as String,
      subscriptionTier: (json['subscription_tier'] as num).toInt(),
    );

Map<String, dynamic> _$UserInfoToJson(UserInfo instance) => <String, dynamic>{
      'user_id': instance.userId,
      'nickname': instance.nickname,
      'avatar_url': instance.avatarUrl,
      'subscription_tier': instance.subscriptionTier,
    };
