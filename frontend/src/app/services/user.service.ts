import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { environment } from '../../environments/environment';

export interface UserProfile {
  id: string;
  name: string;
  email: string;
  createdAt: string;
  updatedAt: string;
}

export interface UserListResponse {
  users: UserProfile[];
  total: number;
  page: number;
  limit: number;
}

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private apiUrl = environment.apiUrl;
  private readonly httpOptions = { withCredentials: true };

  constructor(private http: HttpClient) { }

  private handleError(error: HttpErrorResponse) {
    let errorMessage = 'An error occurred';
    if (error.error instanceof ErrorEvent) {
      errorMessage = `Error: ${error.error.message}`;
    } else {
      errorMessage = `Error Code: ${error.status}\nMessage: ${error.error?.message || error.message}`;
    }
    console.error(errorMessage);
    return throwError(() => new Error(errorMessage));
  }

  getProfile(userId: string): Observable<UserProfile> {
    return this.http.get<UserProfile>(
      `${this.apiUrl}/users/profile/${userId}`,
      this.httpOptions
    ).pipe(
      catchError(this.handleError)
    );
  }

  getCurrentUserProfile(): Observable<UserProfile> {
    return this.http.get<UserProfile>(
      `${this.apiUrl}/users/me`,
      this.httpOptions
    ).pipe(
      catchError(this.handleError)
    );
  }

  listUsers(page: number = 1, limit: number = 10): Observable<UserListResponse> {
    return this.http.get<UserListResponse>(
      `${this.apiUrl}/users/list`,
      {
        ...this.httpOptions,
        params: { page: page.toString(), limit: limit.toString() }
      }
    ).pipe(
      catchError(this.handleError)
    );
  }

  updateProfile(userId: string, profileData: Partial<UserProfile>): Observable<UserProfile> {
    return this.http.put<UserProfile>(
      `${this.apiUrl}/users/profile/${userId}`,
      profileData,
      this.httpOptions
    ).pipe(
      catchError(this.handleError)
    );
  }
}
