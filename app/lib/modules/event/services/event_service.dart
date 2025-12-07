import '../../../core/network/dio_client.dart';
import '../models/event_models.dart';

class EventService {
  final _dio = DioClient.instance;

  Future<List<Event>> getEventTimeline({int page = 1, int pageSize = 20}) async {
    final response = await _dio.get('/api/v1/events/timeline', queryParameters: {
      'page': page,
      'page_size': pageSize,
    });
    return (response.data['list'] as List)
        .map((e) => Event.fromJson(e))
        .toList();
  }

  Future<Event> getEvent(int id) async {
    final response = await _dio.get('/api/v1/events/$id');
    return Event.fromJson(response.data);
  }
}
