import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { login as apiLogin } from '../lib/api';
import { loginSchema, type LoginFormData } from '../lib/schemas';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';

export const LoginPage = () => {
  const { login } = useAuth();
  const navigate = useNavigate();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: '',
      password: '',
    },
  });

  const {
    mutate: loginUser,
    isPending,
    error: backendError,
  } = useMutation({
    mutationFn: async (data: LoginFormData) => {
      return apiLogin(data.email, data.password);
    },
    onSuccess: (user) => {
      login(user);
      navigate('/');
    },
    onError: (error) => {
      console.log(error.message);
    },
  });

  return (
    <div className="flex-1 flex items-center justify-center p-4">
      <Card className="w-full max-w-md shadow-brutalist">
        <CardHeader>
          <CardTitle className="text-2xl font-bold text-center">Login</CardTitle>
        </CardHeader>
        <CardContent>
          {backendError && (
            <div className="bg-destructive/15 text-destructive p-3 rounded-md mb-4 text-sm font-medium border border-destructive">
              {backendError.message}
            </div>
          )}
          <form onSubmit={handleSubmit((data) => loginUser(data))} className="space-y-4" noValidate>
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input id="email" type="email" required {...register('email')} />
              {errors.email && (
                <p className="text-destructive text-sm font-medium">{errors.email.message}</p>
              )}
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input id="password" type="password" required {...register('password')} />
              {errors.password && (
                <p className="text-destructive text-sm font-medium">{errors.password.message}</p>
              )}
            </div>
            <Button
              type="submit"
              disabled={isPending || isSubmitting}
              className="w-full font-bold shadow-brutalist-sm hover:translate-x-[2px] hover:translate-y-[2px] hover:shadow-none transition-all active:translate-x-[4px] active:translate-y-[4px]"
            >
              {isPending || isSubmitting ? 'Logging in...' : 'Login'}
            </Button>
          </form>
        </CardContent>
        <CardFooter className="justify-center">
          <p className="text-muted-foreground text-sm">
            Don't have an account?{' '}
            <Link to="/register" className="text-primary hover:underline font-bold">
              Register
            </Link>
          </p>
        </CardFooter>
      </Card>
    </div>
  );
};
