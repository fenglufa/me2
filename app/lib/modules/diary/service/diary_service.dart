import 'package:dio/dio.dart';
import '../../../core/network/dio_client.dart';
import '../../../core/constants/api_endpoints.dart';
import '../model/diary_stats.dart';
import '../model/diary.dart';

class DiaryService {
  final Dio _dio = DioClient.instance;

  Future<DiaryStats> getDiaryStats() async {
    final response = await _dio.get(ApiEndpoints.diaryStats);
    return DiaryStats.fromJson(response.data);
  }

  Future<DiaryListResponse> getUserDiaries({int page = 1, int pageSize = 20}) async {
    final response = await _dio.get(
      ApiEndpoints.userDiaries,
      queryParameters: {'page': page, 'page_size': pageSize},
    );
    return DiaryListResponse.fromJson(response.data);
  }

  Future<DiaryListResponse> getAvatarDiaries({int page = 1, int pageSize = 20}) async {
    final response = await _dio.get(
      ApiEndpoints.avatarDiaries,
      queryParameters: {'page': page, 'page_size': pageSize},
    );
    return DiaryListResponse.fromJson(response.data);
  }

  Future<Diary> createUserDiary({required String content, String? mood}) async {
    final response = await _dio.post(
      ApiEndpoints.userDiaries,
      data: {
        'content': content,
        if (mood != null && mood.isNotEmpty) 'mood': mood,
      },
    );
    return Diary.fromJson(response.data);
  }
}
