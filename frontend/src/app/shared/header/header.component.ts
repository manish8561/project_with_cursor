import { Component, HostListener } from '@angular/core';
import { RouterLink, RouterLinkActive, Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-header',
  imports: [RouterLink, RouterLinkActive, CommonModule],
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss'
})
export class HeaderComponent {
  isDropdownOpen = false;

  constructor(
    private router: Router,
    private authService: AuthService
  ) { }

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
    const user = this.authService.getUser();
    if (user && user.name) {
      return user.name.split(' ').map((n: string) => n[0]).join('').toUpperCase();
    }
    return 'U';
  }

  getUserName(): string {
    const user = this.authService.getUser();
    return user ? user.name : 'User';
  }

  logout() {
    // Use the auth service to logout
    this.authService.logout();

    // Close dropdown
    this.isDropdownOpen = false;

    // Navigate to login page
    this.router.navigate(['/login']);
  }
}
