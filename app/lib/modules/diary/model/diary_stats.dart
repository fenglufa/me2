import 'package:json_annotation/json_annotation.dart';

part 'diary_stats.g.dart';

@JsonSerializable()
class DiaryStats {
  @JsonKey(name: 'total_count')
  final int totalCount;
  @JsonKey(name: 'avatar_count')
  final int avatarCount;
  @JsonKey(name: 'user_count')
  final int userCount;

  DiaryStats({
    required this.totalCount,
    required this.avatarCount,
    required this.userCount,
  });

  factory DiaryStats.fromJson(Map<String, dynamic> json) =>
      _$DiaryStatsFromJson(json);
  Map<String, dynamic> toJson() => _$DiaryStatsToJson(this);
}
