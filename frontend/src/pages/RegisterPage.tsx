import { Button } from '@/components/ui/button';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
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
            <div className="space-y-2">
              <Label htmlFor="username">Username</Label>
              <Input id="username" {...register('username')} required />
              {errors.username && (
                <p className="text-destructive text-sm font-medium">{errors.username.message}</p>
              )}
            </div>

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
              <Label htmlFor="password">Password</Label>
              <Input id="password" type="password" {...register('password')} required />
              {errors.password && (
                <p className="text-destructive text-sm font-medium">{errors.password.message}</p>
              )}
            </div>

            <div className="flex flex-row items-center space-x-2 pt-2">
              <input
                type="checkbox"
                id="agbAccepted"
                {...register('agbAccepted')}
                className="h-4 w-4 rounded border-gray-300 text-primary focus:ring-primary bg-background shadow-sm"
              />
              <Label htmlFor="agbAccepted" className="text-sm cursor-pointer leading-5">
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
              <p className="text-destructive text-sm font-medium">{errors.agbAccepted.message}</p>
            )}
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
