import { Button } from '@/components/shared/Button';
import { FormInput } from '@/components/shared/FormInput';
import { ErrorMessage } from '@/components/shared/ErrorMessage';
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
import { AuthLayout } from './AuthLayout';
import styles from './AuthLayout.module.css';

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
    <AuthLayout
      title="Register"
      subtitle="Create an account to start trading"
      footer={
        <p className={styles.footerText}>
          Already have an account?{' '}
          <Link to="/login" className={styles.link}>
            Login
          </Link>
        </p>
      }
    >
      {backendError && (
        <ErrorMessage className={styles.errorMessage} message={backendError.message} />
      )}
      <form
        onSubmit={handleSubmit((data) => registerUser(data))}
        className={styles.form}
        noValidate
      >
        <FormInput
          label="Username"
          id="username"
          required
          register={register}
          error={errors.username?.message}
          data-testid="username-input"
        />

        <div className={styles.nameFieldsGrid}>
          <FormInput
            label="First Name"
            id="firstName"
            required
            register={register}
            error={errors.firstName?.message}
            data-testid="first-name-input"
          />
          <FormInput
            label="Last Name"
            id="lastName"
            required
            register={register}
            error={errors.lastName?.message}
            data-testid="last-name-input"
          />
        </div>

        <FormInput
          label="Password"
          id="password"
          type="password"
          required
          register={register}
          error={errors.password?.message}
          data-testid="password-input"
        />

        <FormField error={errors.agbAccepted?.message}>
          <div className={styles.checkboxContainer}>
            <Checkbox id="agbAccepted" {...register('agbAccepted')} data-testid="agb-checkbox" />
            <Label htmlFor="agbAccepted" className={styles.checkboxLabel}>
              I accept the{' '}
              <Link to="/agb" className={styles.link}>
                Terms and Conditions (AGB)
              </Link>{' '}
              and{' '}
              <Link to="/privacy" className={styles.link}>
                Privacy Policy
              </Link>
              .
            </Label>
          </div>
        </FormField>
        <Button
          type="submit"
          className={styles.submitButton}
          disabled={isSubmitting || isRegistering}
          data-testid="register-submit"
        >
          {isSubmitting || isRegistering ? 'Creating Account...' : 'Register'}
        </Button>
      </form>
    </AuthLayout>
  );
};
