import 'package:flutter/material.dart';
import '../models/world_models.dart';
import '../services/world_service.dart';

class WorldController extends ChangeNotifier {
  final _service = WorldService();

  List<WorldMap> maps = [];
  List<WorldRegion> regions = [];
  List<WorldScene> scenes = [];
  bool isLoading = false;
  int? selectedMapId;

  Future<void> loadMaps() async {
    try {
      isLoading = true;
      notifyListeners();
      maps = await _service.getMaps();
      if (maps.isNotEmpty) {
        selectedMapId = maps.first.id;
        await loadRegions(maps.first.id);
      }
    } catch (e) {
      debugPrint('加载世界地图失败: $e');
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }

  Future<void> loadRegions(int mapId) async {
    try {
      isLoading = true;
      notifyListeners();
      selectedMapId = mapId;
      regions = await _service.getRegions(mapId);
    } catch (e) {
      debugPrint('加载区域失败: $e');
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }

  Future<void> loadScenes(int regionId) async {
    try {
      isLoading = true;
      notifyListeners();
      scenes = await _service.getScenes(regionId);
    } catch (e) {
      debugPrint('加载场景失败: $e');
    } finally {
      isLoading = false;
      notifyListeners();
    }
  }
}
