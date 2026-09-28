import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment';

export interface NotificationPreference {
  userId: string;
  emailEnabled: boolean;
  updatedAt: string;
}

export interface NotificationRecord {
  id: string;
  type: string;
  channel: string;
  status: 'queued' | 'sent' | 'failed' | 'skipped';
  createdAt: string;
  updatedAt: string;
  errorSummary?: string;
}

export interface NotificationHistory {
  items: NotificationRecord[];
  total: number;
  page: number;
  size: number;
}

@Injectable({ providedIn: 'root' })
export class NotificationService {
  private readonly apiUrl = `${environment.apiUrl}/notifications`;
  private readonly httpOptions = { withCredentials: true };

  constructor(private readonly http: HttpClient) {}

  getPreference(): Observable<NotificationPreference> {
    return this.http.get<NotificationPreference>(
      `${this.apiUrl}/preferences`,
      this.httpOptions,
    );
  }

  updatePreference(emailEnabled: boolean): Observable<NotificationPreference> {
    return this.http.put<NotificationPreference>(
      `${this.apiUrl}/preferences`,
      { emailEnabled },
      this.httpOptions,
    );
  }

  getHistory(page = 1, size = 10): Observable<NotificationHistory> {
    return this.http.get<NotificationHistory>(`${this.apiUrl}/history`, {
      ...this.httpOptions,
      params: { page: page.toString(), size: size.toString() },
    });
  }
}
