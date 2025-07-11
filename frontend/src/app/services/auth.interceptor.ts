import { HttpEvent, HttpHandlerFn, HttpRequest } from '@angular/common/http';
import { Observable } from 'rxjs';
import { inject } from '@angular/core';
import { AuthService } from './auth.service';

export function authInterceptor(req: HttpRequest<unknown>, next: HttpHandlerFn): Observable<HttpEvent<unknown>> {
    try {
        const authService = inject(AuthService);
        const token = authService.getToken();
        console.log("token", token);

        if (token) {
            const cloned = req.clone({
                headers: req.headers.set('Authorization', `Bearer ${token}`)
            });
            return next(cloned);
        }
    } catch (error) {
        // Optionally log the error for debugging
        console.error('AuthInterceptor error:', error);
    }

    return next(req);
} 