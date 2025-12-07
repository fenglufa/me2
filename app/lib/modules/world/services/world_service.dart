import '../../../core/network/dio_client.dart';
import '../models/world_models.dart';

class WorldService {
  final _dio = DioClient.instance;

  Future<List<WorldMap>> getMaps({int page = 1, int pageSize = 20}) async {
    final response = await _dio.get('/api/v1/world/maps', queryParameters: {
      'page': page,
      'page_size': pageSize,
      'only_active': true,
    });
    return (response.data['list'] as List)
        .map((e) => WorldMap.fromJson(e))
        .toList();
  }

  Future<List<WorldRegion>> getRegions(int mapId,
      {int page = 1, int pageSize = 20}) async {
    final response = await _dio.get('/api/v1/world/regions', queryParameters: {
      'map_id': mapId,
      'page': page,
      'page_size': pageSize,
      'only_active': true,
    });
    return (response.data['list'] as List)
        .map((e) => WorldRegion.fromJson(e))
        .toList();
  }

  Future<List<WorldScene>> getScenes(int regionId,
      {int page = 1, int pageSize = 20}) async {
    final response = await _dio.get('/api/v1/world/scenes', queryParameters: {
      'region_id': regionId,
      'page': page,
      'page_size': pageSize,
      'only_active': true,
    });
    return (response.data['list'] as List)
        .map((e) => WorldScene.fromJson(e))
        .toList();
  }

  Future<WorldRegion> getRegion(int id) async {
    final response = await _dio.get('/api/v1/world/regions/$id');
    return WorldRegion.fromJson(response.data);
  }

  Future<WorldScene> getScene(int id) async {
    final response = await _dio.get('/api/v1/world/scenes/$id');
    return WorldScene.fromJson(response.data);
  }
}
