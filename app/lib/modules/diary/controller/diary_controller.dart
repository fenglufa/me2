import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../service/diary_service.dart';
import '../model/diary_stats.dart';
import '../model/diary.dart';

final diaryServiceProvider = Provider((ref) => DiaryService());

final diaryStatsProvider = FutureProvider<DiaryStats>((ref) async {
  final diaryService = ref.read(diaryServiceProvider);
  return await diaryService.getDiaryStats();
});

final userDiariesProvider = FutureProvider<DiaryListResponse>((ref) async {
  final diaryService = ref.read(diaryServiceProvider);
  return await diaryService.getUserDiaries();
});

final avatarDiariesProvider = FutureProvider<DiaryListResponse>((ref) async {
  final diaryService = ref.read(diaryServiceProvider);
  return await diaryService.getAvatarDiaries();
});
