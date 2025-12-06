import 'package:json_annotation/json_annotation.dart';

part 'diary.g.dart';

@JsonSerializable()
class Diary {
  final int id;
  @JsonKey(name: 'avatar_id')
  final int avatarId;
  final String type;
  final String date;
  final String title;
  final String content;
  final String mood;
  @JsonKey(name: 'created_at')
  final int createdAt;

  Diary({
    required this.id,
    required this.avatarId,
    required this.type,
    required this.date,
    required this.title,
    required this.content,
    required this.mood,
    required this.createdAt,
  });

  factory Diary.fromJson(Map<String, dynamic> json) => _$DiaryFromJson(json);
  Map<String, dynamic> toJson() => _$DiaryToJson(this);
}

@JsonSerializable()
class DiaryListResponse {
  final int total;
  final int page;
  @JsonKey(name: 'page_size')
  final int pageSize;
  final List<Diary> list;

  DiaryListResponse({
    required this.total,
    required this.page,
    required this.pageSize,
    required this.list,
  });

  factory DiaryListResponse.fromJson(Map<String, dynamic> json) =>
      _$DiaryListResponseFromJson(json);
  Map<String, dynamic> toJson() => _$DiaryListResponseToJson(this);
}
