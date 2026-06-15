import { Button } from '@/components/shared/Button';
import { Card } from '@/components/shared/Card';
import { FormInput } from '@/components/shared/FormInput';
import { Label } from '@/components/shared/Label';
import { Checkbox } from '@/components/shared/Checkbox';
import { FormField } from '@/components/shared/FormField';
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
            <FormInput
              label="Username"
              id="username"
              required
              register={register}
              error={errors.username?.message}
            />

            <div className="grid grid-cols-2 gap-4">
              <FormInput
                label="First Name"
                id="firstName"
                required
                register={register}
                error={errors.firstName?.message}
              />
              <FormInput
                label="Last Name"
                id="lastName"
                required
                register={register}
                error={errors.lastName?.message}
              />
            </div>

            <FormInput
              label="Password"
              id="password"
              type="password"
              required
              register={register}
              error={errors.password?.message}
            />

            <FormField error={errors.agbAccepted?.message}>
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
            </FormField>
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
