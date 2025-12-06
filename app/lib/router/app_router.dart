import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../modules/auth/view/login_page.dart';
import '../modules/home/view/home_page.dart';

final appRouter = GoRouter(
  initialLocation: '/login',
  routes: [
    GoRoute(
      path: '/login',
      builder: (context, state) => const LoginPage(),
    ),
    GoRoute(
      path: '/home',
      builder: (context, state) => const HomePage(),
    ),
  ],
);
