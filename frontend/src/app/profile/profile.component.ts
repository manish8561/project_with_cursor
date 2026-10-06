import { Component, OnInit } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { CommonModule } from '@angular/common';
import {
  NotificationHistory,
  NotificationPreference,
  NotificationService,
} from '../services/notification.service';
import { ToastService } from '../services/toast.service';

@Component({
  selector: 'app-profile',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss'],
})
export class ProfileComponent implements OnInit {
  user: any;
  errorMessage: string = '';
  notificationError = '';
  preference?: NotificationPreference;
  history?: NotificationHistory;
  savingPreference = false;

  get historyPageCount(): number {
    return this.history ? Math.ceil(this.history.total / this.history.size) : 0;
  }

  constructor(
    private readonly authService: AuthService,
    private readonly notificationService: NotificationService,
    private readonly toastService: ToastService,
  ) {}

  ngOnInit(): void {
    this.authService.getProfile().subscribe({
      next: (profile) => {
        this.user = profile;
      },
      error: (error) => {
        this.toastService.error(
          error.error?.error || 'Failed to load profile.',
        );
      },
    });

    this.loadNotifications();
  }

  loadNotifications(page = 1): void {
    this.notificationError = '';
    this.notificationService.getPreference().subscribe({
      next: (preference) => (this.preference = preference),
      error: () => {
        this.toastService.error('Could not load email preferences.');
      },
    });
    this.notificationService.getHistory(page).subscribe({
      next: (history) => (this.history = history),
      error: () => {
        this.toastService.error('Could not load notification history.');
      },
    });
  }

  updateEmailPreference(event: Event): void {
    const enabled = (event.target as HTMLInputElement).checked;
    this.savingPreference = true;
    this.notificationError = '';
    this.notificationService.updatePreference(enabled).subscribe({
      next: (preference) => {
        this.preference = preference;
        this.savingPreference = false;
        this.toastService.success('Email preferences updated successfully!');
      },
      error: () => {
        this.toastService.error('Could not save email preferences.');
        this.savingPreference = false;
      },
    });
  }

  changeHistoryPage(page: number): void {
    if (
      !this.history ||
      page < 1 ||
      page > Math.ceil(this.history.total / this.history.size)
    ) {
      return;
    }
    this.loadNotifications(page);
  }
}
