import '../../../core/network/dio_client.dart';
import '../models/action_models.dart';

class ActionService {
  final _dio = DioClient.instance;

  Future<ActionLog?> getLastAction() async {
    final response = await _dio.get('/api/v1/action/last');
    if (response.data['id'] == null) {
      return null;
    }
    return ActionLog.fromJson(response.data);
  }
}
