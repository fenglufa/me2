import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import '../modules/auth/view/login_page.dart';
import '../modules/main/view/main_page.dart';
import '../core/storage/token_storage.dart';

final appRouter = GoRouter(
  initialLocation: '/login',
  redirect: (context, state) async {
    final token = await TokenStorage.getToken();
    final isLoggedIn = token != null && token.isNotEmpty;
    final isLoginPage = state.matchedLocation == '/login';

    if (isLoggedIn && isLoginPage) {
      return '/main';
    }
    if (!isLoggedIn && !isLoginPage) {
      return '/login';
    }
    return null;
  },
  routes: [
    GoRoute(
      path: '/login',
      builder: (context, state) => const LoginPage(),
    ),
    GoRoute(
      path: '/main',
      builder: (context, state) => const MainPage(),
    ),
  ],
);
