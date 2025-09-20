import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';
import { map } from 'rxjs/operators';

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
    user: {
        id: string;
        email: string;
        name: string;
    };
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
                        localStorage.setItem('user', JSON.stringify(response.user));
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
                        localStorage.setItem('user', JSON.stringify(response.user));
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
            localStorage.removeItem('user');
        }
    }

    isLoggedIn(): boolean {
        return this.isBrowser() && !!localStorage.getItem('token');
    }

    getToken(): string | null {
        return this.isBrowser() ? localStorage.getItem('token') : null;
    }

    getUser(): any {
        if (!this.isBrowser()) return null;
        const user = localStorage.getItem('user');
        return user ? JSON.parse(user) : null;
    }

    getProfile(): Observable<any> {
        return this.http.get<any>(`${this.apiUrl}/users/profile`);
    }

    isBrowser(): boolean {
        return typeof window !== 'undefined' && !!window.localStorage;
    }
} 