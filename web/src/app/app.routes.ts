import { Routes } from '@angular/router';
import { Dashboard } from './features/dashboard/dashboard';
import { Home } from './features/home/home';

export const routes: Routes = [
  {
    path: '',
    component: Dashboard,
  },
  {
    path: 'hello',
    component: Home,
  },
];
