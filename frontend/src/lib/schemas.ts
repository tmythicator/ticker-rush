import { z } from 'zod';

export const registerSchema = z.object({
  email: z.email('Invalid email adress'),
  password: z.string().min(8, 'Password must be at least 8 characters long'),
  firstName: z.string().min(1, 'First name is required'),
  lastName: z.string().min(1, 'Last name is required'),
});

export const loginSchema = z.object({
  email: z.email('Invalid email adress'),
  password: z.string().min(8, 'Password must be at least 8 characters long'),
});

export type RegisterFormData = z.infer<typeof registerSchema>;
export type LoginFormData = z.infer<typeof loginSchema>;
