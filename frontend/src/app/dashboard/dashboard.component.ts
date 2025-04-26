import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [RouterLink],
  template: `
    <div class="dashboard-container">
      <h1>Welcome to Dashboard</h1>
      <p>You have successfully logged in!</p>
      <a routerLink="/" class="btn btn-primary">Back to Home</a>
    </div>
  `,
  styles: [`
    .dashboard-container {
      max-width: 800px;
      margin: 2rem auto;
      padding: 2rem;
      text-align: center;
    }
  `]
})
export class DashboardComponent {}
