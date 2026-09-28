import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of } from 'rxjs';

import { ProfileComponent } from './profile.component';
import { AuthService } from '../services/auth.service';
import { NotificationService } from '../services/notification.service';

describe('ProfileComponent', () => {
  let component: ProfileComponent;
  let fixture: ComponentFixture<ProfileComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ProfileComponent],
      providers: [
        { provide: AuthService, useValue: { getProfile: () => of({}) } },
        {
          provide: NotificationService,
          useValue: {
            getPreference: () =>
              of({ userId: 'user-1', emailEnabled: true, updatedAt: '' }),
            getHistory: () => of({ items: [], total: 0, page: 1, size: 10 }),
            updatePreference: () =>
              of({ userId: 'user-1', emailEnabled: true, updatedAt: '' }),
          },
        },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(ProfileComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
