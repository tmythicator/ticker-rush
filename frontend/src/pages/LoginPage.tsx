import { Button } from '@/components/shared/Button';
import { Card } from '@/components/shared/Card';
import { FormInput } from '@/components/shared/FormInput';
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
    <div className="flex flex-1 items-center justify-center p-4">
      <Card className="shadow-brutalist w-full max-w-md p-6">
        <div className="mb-6 flex flex-col space-y-1.5 text-center">
          <h2 className="text-2xl font-bold">Login</h2>
        </div>
        <div>
          {backendError && (
            <div className="mb-4 rounded-md border border-destructive bg-destructive/15 p-3 text-sm font-medium text-destructive">
              {backendError.message}
            </div>
          )}
          <form onSubmit={handleSubmit((data) => loginUser(data))} className="space-y-4" noValidate>
            <FormInput
              label="Username"
              id="username"
              type="text"
              required
              register={register}
              error={errors.username?.message}
            />
            <FormInput
              label="Password"
              id="password"
              type="password"
              required
              register={register}
              error={errors.password?.message}
            />
            <Button type="submit" disabled={isPending || isSubmitting} className="w-full">
              {isPending || isSubmitting ? 'Logging in...' : 'Login'}
            </Button>
          </form>
        </div>
        <div className="mt-6 flex justify-center border-t border-border/50 pt-4">
          <p className="text-sm text-muted-foreground">
            Don't have an account?{' '}
            <Link to="/register" className="font-bold text-primary hover:underline">
              Register
            </Link>
          </p>
        </div>
      </Card>
    </div>
  );
};
