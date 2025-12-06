// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'diary.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Diary _$DiaryFromJson(Map<String, dynamic> json) => Diary(
      id: (json['id'] as num).toInt(),
      avatarId: (json['avatar_id'] as num).toInt(),
      type: json['type'] as String,
      date: json['date'] as String,
      title: json['title'] as String,
      content: json['content'] as String,
      mood: json['mood'] as String,
      replyContent: json['reply_content'] as String?,
      createdAt: (json['created_at'] as num).toInt(),
    );

Map<String, dynamic> _$DiaryToJson(Diary instance) => <String, dynamic>{
      'id': instance.id,
      'avatar_id': instance.avatarId,
      'type': instance.type,
      'date': instance.date,
      'title': instance.title,
      'content': instance.content,
      'mood': instance.mood,
      'reply_content': instance.replyContent,
      'created_at': instance.createdAt,
    };

DiaryListResponse _$DiaryListResponseFromJson(Map<String, dynamic> json) =>
    DiaryListResponse(
      total: (json['total'] as num).toInt(),
      page: (json['page'] as num).toInt(),
      pageSize: (json['page_size'] as num).toInt(),
      list: (json['list'] as List<dynamic>)
          .map((e) => Diary.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$DiaryListResponseToJson(DiaryListResponse instance) =>
    <String, dynamic>{
      'total': instance.total,
      'page': instance.page,
      'page_size': instance.pageSize,
      'list': instance.list,
    };
