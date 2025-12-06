// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'diary_stats.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

DiaryStats _$DiaryStatsFromJson(Map<String, dynamic> json) => DiaryStats(
      totalCount: (json['total_count'] as num).toInt(),
      avatarCount: (json['avatar_count'] as num).toInt(),
      userCount: (json['user_count'] as num).toInt(),
    );

Map<String, dynamic> _$DiaryStatsToJson(DiaryStats instance) =>
    <String, dynamic>{
      'total_count': instance.totalCount,
      'avatar_count': instance.avatarCount,
      'user_count': instance.userCount,
    };
