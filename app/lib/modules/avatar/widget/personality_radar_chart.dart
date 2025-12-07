import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';

class PersonalityRadarChart extends StatelessWidget {
  final double warmth;
  final double adventurous;
  final double social;
  final double creative;
  final double calm;
  final double energetic;

  const PersonalityRadarChart({
    super.key,
    required this.warmth,
    required this.adventurous,
    required this.social,
    required this.creative,
    required this.calm,
    required this.energetic,
  });

  @override
  Widget build(BuildContext context) {
    return RadarChart(
      RadarChartData(
        radarShape: RadarShape.polygon,
        tickCount: 5,
        ticksTextStyle: const TextStyle(fontSize: 10, color: Colors.transparent),
        tickBorderData: const BorderSide(color: Colors.grey, width: 1),
        gridBorderData: const BorderSide(color: Colors.grey, width: 1),
        radarBorderData: const BorderSide(color: Colors.purple, width: 2),
        radarBackgroundColor: Colors.transparent,
        radarTouchData: RadarTouchData(enabled: false),
        dataSets: [
          RadarDataSet(
            fillColor: Colors.transparent,
            borderColor: Colors.transparent,
            dataEntries: List.filled(6, const RadarEntry(value: 1.0)),
          ),
          RadarDataSet(
            fillColor: Colors.purple.withValues(alpha: 0.2),
            borderColor: Colors.purple,
            borderWidth: 2,
            dataEntries: [
              RadarEntry(value: warmth),
              RadarEntry(value: adventurous),
              RadarEntry(value: social),
              RadarEntry(value: creative),
              RadarEntry(value: calm),
              RadarEntry(value: energetic),
            ],
          ),
        ],
        getTitle: (index, angle) {
          const labels = ['温暖度', '冒险性', '社交性', '创造力', '平静度', '活力值'];
          return RadarChartTitle(text: labels[index]);
        },
        titleTextStyle: const TextStyle(fontSize: 12, color: Colors.black87),
      ),
    );
  }
}
