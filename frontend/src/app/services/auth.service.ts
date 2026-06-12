import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';
import { catchError, map } from 'rxjs/operators';

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
    token: string;
    message?: string;
}

@Injectable({
    providedIn: 'root'
})
export class AuthService {
    private apiUrl = environment.apiUrl;

    constructor(private http: HttpClient) { }

    login(credentials: LoginRequest): Observable<AuthResponse> {
        return this.http.post<AuthResponse>(`${this.apiUrl}/auth/login`, credentials).pipe(
            map(response => {
                if (response.token) {
                    if (this.isBrowser()) {
                        localStorage.setItem('token', response.token);
                    }
                    return response;
                }
                throw new Error('Login failed');
            })
        );
    }

    register(userData: RegisterRequest): Observable<AuthResponse> {
        return this.http.post<AuthResponse>(`${this.apiUrl}/auth/register`, userData).pipe(
            map(response => {
                if (response.token) {
                    if (this.isBrowser()) {
                        localStorage.setItem('token', response.token);
                    }
                    return response;
                }
                throw new Error('Registration failed');
            })
        );
    }

    logout(): void {
        if (this.isBrowser()) {
            localStorage.removeItem('token');
        }
    }

    isLoggedIn(): boolean {
        return this.isBrowser() && !!localStorage.getItem('token');
    }

    getToken(): string | null {
        return this.isBrowser() ? localStorage.getItem('token') : null;
    }

    getUserId(): string | null {
        const token = this.getToken();
        if (!token) {
            return null;
        }

        try {
            const payload = token.split('.')[1];
            const decoded = JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/')));
            return decoded.user_id ?? null;
        } catch {
            return null;
        }
    }

    getProfile(): Observable<any> {
        const userId = this.getUserId();
        if (!userId) {
            throw new Error('User not authenticated');
        }

        return this.http.get<any>(`${this.apiUrl}/users/profile/${userId}`, {
            headers: {
                'Authorization': `Bearer ${this.getToken()}`
            }
        }).pipe(
            catchError(error => {
                console.error('Error fetching profile:', error);
                throw error;
            })
        );
    }

    isBrowser(): boolean {
        return typeof window !== 'undefined' && !!window.localStorage;
    }
} 