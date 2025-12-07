class WorldMap {
  final int id;
  final String name;
  final String description;
  final String coverImage;
  final bool isActive;
  final int createdAt;

  WorldMap({
    required this.id,
    required this.name,
    required this.description,
    required this.coverImage,
    required this.isActive,
    required this.createdAt,
  });

  factory WorldMap.fromJson(Map<String, dynamic> json) => WorldMap(
        id: json['id'],
        name: json['name'],
        description: json['description'],
        coverImage: json['cover_image'],
        isActive: json['is_active'],
        createdAt: json['created_at'],
      );
}

class WorldRegion {
  final int id;
  final int mapId;
  final String name;
  final String description;
  final String coverImage;
  final String atmosphere;
  final List<String> tags;
  final bool isActive;
  final int createdAt;

  WorldRegion({
    required this.id,
    required this.mapId,
    required this.name,
    required this.description,
    required this.coverImage,
    required this.atmosphere,
    required this.tags,
    required this.isActive,
    required this.createdAt,
  });

  factory WorldRegion.fromJson(Map<String, dynamic> json) => WorldRegion(
        id: json['id'],
        mapId: json['map_id'],
        name: json['name'],
        description: json['description'],
        coverImage: json['cover_image'],
        atmosphere: json['atmosphere'],
        tags: List<String>.from(json['tags'] ?? []),
        isActive: json['is_active'],
        createdAt: json['created_at'],
      );
}

class SceneFeatures {
  final bool hasWifi;
  final bool hasFood;
  final bool hasSeating;
  final bool isIndoor;
  final bool isQuiet;
  final int comfortLevel;
  final int socialLevel;

  SceneFeatures({
    required this.hasWifi,
    required this.hasFood,
    required this.hasSeating,
    required this.isIndoor,
    required this.isQuiet,
    required this.comfortLevel,
    required this.socialLevel,
  });

  factory SceneFeatures.fromJson(Map<String, dynamic> json) => SceneFeatures(
        hasWifi: json['has_wifi'],
        hasFood: json['has_food'],
        hasSeating: json['has_seating'],
        isIndoor: json['is_indoor'],
        isQuiet: json['is_quiet'],
        comfortLevel: json['comfort_level'],
        socialLevel: json['social_level'],
      );
}

class WorldScene {
  final int id;
  final int regionId;
  final String name;
  final String description;
  final String coverImage;
  final String atmosphere;
  final List<String> tags;
  final List<String> suitableActions;
  final int capacity;
  final bool isActive;
  final SceneFeatures features;
  final int createdAt;

  WorldScene({
    required this.id,
    required this.regionId,
    required this.name,
    required this.description,
    required this.coverImage,
    required this.atmosphere,
    required this.tags,
    required this.suitableActions,
    required this.capacity,
    required this.isActive,
    required this.features,
    required this.createdAt,
  });

  factory WorldScene.fromJson(Map<String, dynamic> json) => WorldScene(
        id: json['id'],
        regionId: json['region_id'],
        name: json['name'],
        description: json['description'],
        coverImage: json['cover_image'],
        atmosphere: json['atmosphere'],
        tags: List<String>.from(json['tags'] ?? []),
        suitableActions: List<String>.from(json['suitable_actions'] ?? []),
        capacity: json['capacity'],
        isActive: json['is_active'],
        features: SceneFeatures.fromJson(json['features']),
        createdAt: json['created_at'],
      );
}
