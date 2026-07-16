import { Button } from '@/components/shared/Button';
import { FormInput } from '@/components/shared/FormInput';
import { ErrorMessage } from '@/components/shared/ErrorMessage';
import { useAuth } from '@/hooks/useAuth';
import { login as apiLogin } from '@/lib/api';
import { loginSchema, type LoginFormData } from '@/lib/schemas';
import { zodResolver } from '@hookform/resolvers/zod';
import { useMutation } from '@tanstack/react-query';
import { useForm } from 'react-hook-form';
import { Link, useNavigate } from 'react-router-dom';
import { AuthLayout } from './AuthLayout';
import styles from './AuthLayout.module.css';

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
  });

  return (
    <AuthLayout
      title="Login"
      footer={
        <p className={styles.footerText}>
          Don't have an account?{' '}
          <Link to="/register" className={styles.link}>
            Register
          </Link>
        </p>
      }
    >
      {backendError && (
        <ErrorMessage className={styles.errorMessage} message={backendError.message} />
      )}
      <form onSubmit={handleSubmit((data) => loginUser(data))} className={styles.form} noValidate>
        <FormInput
          label="Username"
          id="username"
          type="text"
          required
          register={register}
          error={errors.username?.message}
          data-testid="username-input"
        />
        <FormInput
          label="Password"
          id="password"
          type="password"
          required
          register={register}
          error={errors.password?.message}
          data-testid="password-input"
        />
        <Button
          type="submit"
          disabled={isPending || isSubmitting}
          className={styles.submitButton}
          data-testid="login-submit"
        >
          {isPending || isSubmitting ? 'Logging in...' : 'Login'}
        </Button>
      </form>
    </AuthLayout>
  );
};
