import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../model/diary.dart';
import '../service/diary_service.dart';

final diaryDetailProvider = FutureProvider.family<Diary, int>((ref, diaryId) async {
  final diaryService = DiaryService();
  return await diaryService.getDiary(diaryId);
});

class DiaryDetailPage extends ConsumerWidget {
  final int diaryId;

  const DiaryDetailPage({super.key, required this.diaryId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final diaryAsync = ref.watch(diaryDetailProvider(diaryId));

    return Scaffold(
      appBar: AppBar(
        title: const Text('日记详情'),
      ),
      body: diaryAsync.when(
        data: (diary) => SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    diary.title.isEmpty ? '无标题' : diary.title,
                    style: const TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  if (diary.mood.isNotEmpty)
                    Chip(label: Text(diary.mood)),
                ],
              ),
              const SizedBox(height: 8),
              Text(
                diary.date,
                style: TextStyle(
                  fontSize: 14,
                  color: Colors.grey.shade600,
                ),
              ),
              const SizedBox(height: 24),
              Text(
                diary.content,
                style: const TextStyle(fontSize: 16, height: 1.5),
              ),
              if (diary.replyContent != null && diary.replyContent!.isNotEmpty) ...[
                const SizedBox(height: 32),
                const Divider(),
                const SizedBox(height: 16),
                Row(
                  children: [
                    Icon(Icons.reply, color: Colors.blue.shade400),
                    const SizedBox(width: 8),
                    const Text(
                      '分身回复',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 12),
                Container(
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Colors.blue.shade50,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    diary.replyContent!,
                    style: const TextStyle(fontSize: 16, height: 1.5),
                  ),
                ),
              ],
            ],
          ),
        ),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, _) => Center(child: Text('加载失败: $error')),
      ),
    );
  }
}
