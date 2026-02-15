import { useNavigate, Link } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation } from '@tanstack/react-query';
import { register as apiRegister, login as apiLogin } from '@/lib/api';
import { useAuth } from '@/hooks/useAuth';
import { type RegisterFormData, registerSchema } from '@/lib/schemas';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';

export const RegisterPage = () => {
  const { login } = useAuth();
  const navigate = useNavigate();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      email: '',
      password: '',
      firstName: '',
      lastName: '',
    },
  });

  const {
    mutate: registerUser,
    isPending: isRegistering,
    error: backendError,
  } = useMutation({
    mutationFn: async (data: RegisterFormData) => {
      await apiRegister(data.email, data.password, data.firstName, data.lastName);
      return apiLogin(data.email, data.password);
    },
    onSuccess: (user) => {
      login(user);
      navigate('/');
    },
    onError: (error) => {
      console.error(error.message);
    },
  });

  return (
    <div className="flex-1 flex items-center justify-center p-4">
      <Card className="w-full max-w-md shadow-brutalist">
        <CardHeader>
          <CardTitle className="text-2xl font-bold text-center">Register</CardTitle>
          <CardDescription className="text-center">
            Create an account to start trading
          </CardDescription>
        </CardHeader>
        <CardContent>
          {backendError && (
            <div className="bg-destructive/15 text-destructive p-3 rounded-md mb-4 text-sm font-medium border border-destructive">
              {backendError.message}
            </div>
          )}
          <form
            onSubmit={handleSubmit((data) => registerUser(data))}
            className="space-y-4"
            noValidate
          >
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="firstName">First Name</Label>
                <Input id="firstName" {...register('firstName')} required />
                {errors.firstName && (
                  <p className="text-destructive text-sm font-medium">{errors.firstName.message}</p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="lastName">Last Name</Label>
                <Input id="lastName" {...register('lastName')} required />
                {errors.lastName && (
                  <p className="text-destructive text-sm font-medium">{errors.lastName.message}</p>
                )}
              </div>
            </div>
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input id="email" type="email" {...register('email')} required />
              {errors.email && (
                <p className="text-destructive text-sm font-medium">{errors.email.message}</p>
              )}
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input id="password" type="password" {...register('password')} required />
              {errors.password && (
                <p className="text-destructive text-sm font-medium">{errors.password.message}</p>
              )}
            </div>
            <Button
              type="submit"
              className="w-full font-bold shadow-brutalist-sm hover:translate-x-[2px] hover:translate-y-[2px] hover:shadow-none transition-all active:translate-x-[4px] active:translate-y-[4px]"
              disabled={isSubmitting || isRegistering}
            >
              {isSubmitting || isRegistering ? 'Creating Account...' : 'Register'}
            </Button>
          </form>
        </CardContent>
        <CardFooter className="justify-center">
          <p className="text-muted-foreground text-sm">
            Already have an account?{' '}
            <Link to="/login" className="text-primary hover:underline font-bold">
              Login
            </Link>
          </p>
        </CardFooter>
      </Card>
    </div>
  );
};
