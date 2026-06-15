import { Button } from '@/components/shared/Button';
import { Card } from '@/components/shared/Card';
import { Input } from '@/components/shared/Input';
import { Label } from '@/components/shared/Label';
import { Checkbox } from '@/components/shared/Checkbox';
import { useAuth } from '@/hooks/useAuth';
import { login as apiLogin, register as apiRegister } from '@/lib/api';
import { registerSchema, type RegisterFormData } from '@/lib/schemas';
import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { Link, useNavigate } from 'react-router-dom';

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
      username: '',
      password: '',
      firstName: '',
      lastName: '',
      agbAccepted: false,
    },
  });

  const {
    mutate: registerUser,
    isPending: isRegistering,
    error: backendError,
  } = useMutation({
    mutationFn: async (data: RegisterFormData) => {
      await apiRegister({
        username: data.username,
        password: data.password,
        first_name: data.firstName,
        last_name: data.lastName,
        agb_accepted: data.agbAccepted,
        website: '',
      });
      return apiLogin({ username: data.username, password: data.password });
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
    <div className="flex flex-1 items-center justify-center p-4">
      <Card className="shadow-brutalist w-full max-w-md p-6">
        <div className="mb-6 flex flex-col space-y-1.5 text-center">
          <h2 className="text-2xl font-bold">Register</h2>
          <p className="text-sm text-muted-foreground">Create an account to start trading</p>
        </div>
        <div>
          {backendError && (
            <div className="mb-4 rounded-md border border-destructive bg-destructive/15 p-3 text-sm font-medium text-destructive">
              {backendError.message}
            </div>
          )}
          <form
            onSubmit={handleSubmit((data) => registerUser(data))}
            className="space-y-4"
            noValidate
          >
            <div className="space-y-2">
              <Label htmlFor="username">Username</Label>
              <Input id="username" {...register('username')} required />
              {errors.username && (
                <p className="text-sm font-medium text-destructive">{errors.username.message}</p>
              )}
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="firstName">First Name</Label>
                <Input id="firstName" {...register('firstName')} required />
                {errors.firstName && (
                  <p className="text-sm font-medium text-destructive">{errors.firstName.message}</p>
                )}
              </div>
              <div className="space-y-2">
                <Label htmlFor="lastName">Last Name</Label>
                <Input id="lastName" {...register('lastName')} required />
                {errors.lastName && (
                  <p className="text-sm font-medium text-destructive">{errors.lastName.message}</p>
                )}
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input id="password" type="password" {...register('password')} required />
              {errors.password && (
                <p className="text-sm font-medium text-destructive">{errors.password.message}</p>
              )}
            </div>

            <div className="flex flex-row items-center space-x-2 pt-2">
              <Checkbox id="agbAccepted" {...register('agbAccepted')} className="h-4 w-4" />
              <Label htmlFor="agbAccepted" className="cursor-pointer text-sm leading-5">
                I accept the{' '}
                <Link to="/agb" className="text-primary hover:underline">
                  Terms and Conditions (AGB)
                </Link>{' '}
                and{' '}
                <Link to="/privacy" className="text-primary hover:underline">
                  Privacy Policy
                </Link>
                .
              </Label>
            </div>
            {errors.agbAccepted && (
              <p className="text-sm font-medium text-destructive">{errors.agbAccepted.message}</p>
            )}
            <Button type="submit" className="w-full" disabled={isSubmitting || isRegistering}>
              {isSubmitting || isRegistering ? 'Creating Account...' : 'Register'}
            </Button>
          </form>
        </div>
        <div className="mt-6 flex justify-center border-t border-border/50 pt-4">
          <p className="text-sm text-muted-foreground">
            Already have an account?{' '}
            <Link to="/login" className="font-bold text-primary hover:underline">
              Login
            </Link>
          </p>
        </div>
      </Card>
    </div>
  );
};
