import 'dart:io';
import 'package:dio/dio.dart';
import 'package:web_socket_channel/io.dart';
import '../../../core/network/dio_client.dart';
import '../../../core/storage/token_storage.dart';
import '../models/dialogue_models.dart';

class DialogueService {
  final Dio _dio = DioClient.instance;

  Future<int> createSession(int avatarId, {String? title}) async {
    final response = await _dio.post(
      '/api/v1/dialogue/sessions',
      data: {
        'avatar_id': avatarId,
        if (title != null) 'title': title,
      },
    );
    return response.data['session_id'];
  }

  Future<List<Session>> getSessions(int avatarId, {int page = 1, int pageSize = 20}) async {
    final response = await _dio.get(
      '/api/v1/dialogue/sessions',
      queryParameters: {
        'avatar_id': avatarId,
        'page': page,
        'page_size': pageSize,
      },
    );
    return (response.data['sessions'] as List).map((s) => Session.fromJson(s)).toList();
  }

  Future<List<Message>> getMessages(int sessionId, {int page = 1, int pageSize = 50}) async {
    final response = await _dio.get(
      '/api/v1/dialogue/messages/$sessionId',
      queryParameters: {
        'page': page,
        'page_size': pageSize,
      },
    );
    return (response.data['messages'] as List).map((m) => Message.fromJson(m)).toList();
  }

  Future<void> deleteSession(int sessionId) async {
    await _dio.delete('/api/v1/dialogue/sessions/$sessionId');
  }

  Future<IOWebSocketChannel> connectStream(int avatarId) async {
    final token = await TokenStorage.getToken();
    final socket = await WebSocket.connect(
      'ws://localhost:8888/api/v1/dialogue/stream?avatar_id=$avatarId',
      headers: {'Authorization': 'Bearer $token'},
    );
    return IOWebSocketChannel(socket);
  }
}
