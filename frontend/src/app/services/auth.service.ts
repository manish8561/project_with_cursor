import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { environment } from '../../environments/environment';
import { catchError, map, tap } from 'rxjs/operators';

export interface LoginRequest {
    email: string;
    password: string;
}

export interface RegisterRequest {
    name: string;
    email: string;
    password: string;
    confirmPassword: string;
}

export interface AuthResponse {
    status: string;
    message?: string;
    user_id?: string;
}

@Injectable({
    providedIn: 'root'
})
export class AuthService {
    private apiUrl = environment.apiUrl;
    private authenticated = false;
    private userId: string | null = null;
    private readonly httpOptions = { withCredentials: true };

    constructor(private http: HttpClient) { }

    login(credentials: LoginRequest): Observable<AuthResponse> {
        return this.http.post<AuthResponse>(`${this.apiUrl}/auth/login`, credentials, this.httpOptions).pipe(
            tap(response => {
                if (response.status === 'success') {
                    this.authenticated = true;
                }
            }),
            map(response => {
                if (response.status === 'success') {
                    return response;
                }
                throw new Error('Login failed');
            })
        );
    }

    register(userData: RegisterRequest): Observable<AuthResponse> {
        return this.http.post<AuthResponse>(`${this.apiUrl}/auth/register`, userData, this.httpOptions).pipe(
            tap(response => {
                if (response.status === 'success') {
                    this.authenticated = true;
                }
            }),
            map(response => {
                if (response.status === 'success') {
                    return response;
                }
                throw new Error('Registration failed');
            })
        );
    }

    logout(): Observable<AuthResponse> {
        return this.http.post<AuthResponse>(`${this.apiUrl}/auth/logout`, {}, this.httpOptions).pipe(
            tap(() => this.clearSession()),
            catchError(error => {
                this.clearSession();
                throw error;
            })
        );
    }

    /** Validates the HttpOnly session cookie with the auth service. */
    checkAuth(): Observable<boolean> {
        return this.http.get<AuthResponse>(`${this.apiUrl}/auth/me`, this.httpOptions).pipe(
            map(response => {
                const ok = response.status === 'success' && !!response.user_id;
                this.authenticated = ok;
                this.userId = ok ? (response.user_id ?? null) : null;
                return ok;
            }),
            catchError(() => {
                this.clearSession();
                return of(false);
            })
        );
    }

    isLoggedIn(): boolean {
        return this.authenticated;
    }

    getUserId(): string | null {
        return this.userId;
    }

    getProfile(): Observable<any> {
        return this.http.get<any>(`${this.apiUrl}/users/me`, this.httpOptions).pipe(
            tap(profile => {
                if (profile?.id) {
                    this.userId = profile.id;
                    this.authenticated = true;
                }
            }),
            catchError(error => {
                console.error('Error fetching profile:', error);
                throw error;
            })
        );
    }

    clearSession(): void {
        this.authenticated = false;
        this.userId = null;
    }

    isBrowser(): boolean {
        return typeof window !== 'undefined';
    }
}
