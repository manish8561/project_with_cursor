import { Injectable } from '@angular/core';
import { ToastrService } from 'ngx-toastr';

@Injectable({
  providedIn: 'root',
})
export class ToastService {
  constructor(private toastr: ToastrService) {}

  success(message: string, title?: string): void {
    this.toastr.success(message, title || 'Success', {
      timeOut: 3000,
      positionClass: 'toast-top-right',
      closeButton: true,
      progressBar: true,
    });
  }

  error(message: string, title?: string): void {
    this.toastr.error(message, title || 'Error', {
      timeOut: 5000,
      positionClass: 'toast-top-right',
      closeButton: true,
      progressBar: true,
    });
  }

  warning(message: string, title?: string): void {
    this.toastr.warning(message, title || 'Warning', {
      timeOut: 4000,
      positionClass: 'toast-top-right',
      closeButton: true,
      progressBar: true,
    });
  }

  info(message: string, title?: string): void {
    this.toastr.info(message, title || 'Info', {
      timeOut: 3000,
      positionClass: 'toast-top-right',
      closeButton: true,
      progressBar: true,
    });
  }

  clear(): void {
    this.toastr.clear();
  }
}
