import { Button } from '@/components/shared/Button';
import { Card } from '@/components/shared/Card';
import { Input } from '@/components/shared/Input';
import { Label } from '@/components/shared/Label';
import { useAuth } from '@/hooks/useAuth';
import { login as apiLogin } from '@/lib/api';
import { loginSchema, type LoginFormData } from '@/lib/schemas';
import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { Link, useNavigate } from 'react-router-dom';

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
      username: '',
      password: '',
    },
  });

  const {
    mutate: loginUser,
    isPending,
    error: backendError,
  } = useMutation({
    mutationFn: async (data: LoginFormData) => {
      const { username, password } = data;
      return apiLogin({ username, password });
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
      <Card className="w-full max-w-md p-6 shadow-brutalist">
        <div className="flex flex-col space-y-1.5 text-center mb-6">
          <h2 className="text-2xl font-bold">Login</h2>
        </div>
        <div>
          {backendError && (
            <div className="bg-destructive/15 text-destructive p-3 rounded-md mb-4 text-sm font-medium border border-destructive">
              {backendError.message}
            </div>
          )}
          <form onSubmit={handleSubmit((data) => loginUser(data))} className="space-y-4" noValidate>
            <div className="space-y-2">
              <Label htmlFor="username">Username</Label>
              <Input id="username" type="text" required {...register('username')} />
              {errors.username && (
                <p className="text-destructive text-sm font-medium">{errors.username.message}</p>
              )}
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input id="password" type="password" required {...register('password')} />
              {errors.password && (
                <p className="text-destructive text-sm font-medium">{errors.password.message}</p>
              )}
            </div>
            <Button type="submit" disabled={isPending || isSubmitting} className="w-full">
              {isPending || isSubmitting ? 'Logging in...' : 'Login'}
            </Button>
          </form>
        </div>
        <div className="flex justify-center mt-6 pt-4 border-t border-border/50">
          <p className="text-muted-foreground text-sm">
            Don't have an account?{' '}
            <Link to="/register" className="text-primary hover:underline font-bold">
              Register
            </Link>
          </p>
        </div>
      </Card>
    </div>
  );
};
