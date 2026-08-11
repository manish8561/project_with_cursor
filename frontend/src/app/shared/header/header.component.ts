import { Component, HostListener, OnInit } from '@angular/core';
import { RouterLink, RouterLinkActive, Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-header',
  imports: [RouterLink, RouterLinkActive, CommonModule],
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss'
})
export class HeaderComponent implements OnInit {
  isDropdownOpen = false;
  profile: { name?: string } | null = null;

  constructor(
    private router: Router,
    private authService: AuthService
  ) { }

  ngOnInit(): void {
    this.authService.checkAuth().subscribe(isLoggedIn => {
      if (!isLoggedIn) {
        return;
      }

      this.authService.getProfile().subscribe({
        next: (profile) => {
          this.profile = profile;
        },
        error: () => {
          this.profile = null;
        }
      });
    });
  }

  @HostListener('document:click', ['$event'])
  onDocumentClick(event: Event) {
    const target = event.target as HTMLElement;
    if (!target.closest('.profile-dropdown')) {
      this.isDropdownOpen = false;
    }
  }

  toggleDropdown() {
    this.isDropdownOpen = !this.isDropdownOpen;
  }

  getUserInitials(): string {
    if (this.profile?.name) {
      return this.profile.name.split(' ').map((n: string) => n[0]).join('').toUpperCase();
    }
    return 'U';
  }

  getUserName(): string {
    return this.profile?.name ?? 'User';
  }

  logout() {
    this.authService.logout().subscribe({
      next: () => {
        this.isDropdownOpen = false;
        this.router.navigate(['/login']);
      },
      error: () => {
        this.isDropdownOpen = false;
        this.router.navigate(['/login']);
      }
    });
  }
}
