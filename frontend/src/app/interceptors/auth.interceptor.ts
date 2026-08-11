import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor,
  HttpErrorResponse
} from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  constructor(
    private router: Router,
    private authService: AuthService
  ) {}

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    request = request.clone({ withCredentials: true });

    return next.handle(request).pipe(
      catchError((error: HttpErrorResponse) => {
        const isAuthProbe =
          request.url.includes('/auth/login') ||
          request.url.includes('/auth/register') ||
          request.url.includes('/auth/me') ||
          request.url.includes('/auth/logout');

        if (error.status === 401 && !isAuthProbe) {
          this.authService.clearSession();
          this.router.navigate(['/login']);
        }
        return throwError(() => error);
      })
    );
  }
}
