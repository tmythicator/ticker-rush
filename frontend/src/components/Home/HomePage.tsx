import { buttonVariants } from '@/components/shared/buttonVariants';
import { IconArrowRight } from '@icons/CustomIcons';
import { Link, Navigate } from 'react-router-dom';
import { HomeChart } from './HomeChart';
import { cn } from '@/lib/utils';
import { useAuth } from '@/hooks/useAuth';

export const HomePage = () => {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return null;
  }

  if (user) {
    return <Navigate to="/profile" replace />;
  }

  return (
    <div className="flex flex-1 flex-col">
      <section className="relative flex flex-1 flex-col items-center justify-center overflow-hidden pb-20 pt-12">
        <div className="container relative z-20 flex flex-col items-center px-4 text-center md:px-6">
          <div className="max-w-3xl space-y-4">
            <h1 className="text-4xl font-black tracking-tighter text-gray-900 dark:bg-gradient-to-r dark:from-white dark:to-gray-500 dark:bg-clip-text dark:text-transparent sm:text-5xl md:text-6xl">
              Play the Markets.
              <br />
              <span className="text-primary">Compete in Monthly Ladders.</span>
            </h1>
            <p className="mx-auto max-w-[700px] text-gray-500 dark:text-gray-400 md:text-xl">
              Experience the thrill of real-time trading{' '}
              <strong className="font-semibold text-gray-900 dark:text-gray-100">risk-free</strong>.
              Compete in{' '}
              <strong className="font-semibold text-gray-900 dark:text-gray-100">
                monthly ladders
              </strong>{' '}
              and prove your skills to the world.
            </p>
          </div>
          <div className="mt-8 flex flex-col gap-2 min-[400px]:flex-row">
            <Link
              to="/register"
              className={cn(buttonVariants({ size: 'lg' }), 'h-12 px-8 text-base')}
            >
              Get Started <IconArrowRight className="ml-2 h-4 w-4" />
            </Link>
            <Link
              to="/leaderboard"
              className={cn(
                buttonVariants({ variant: 'outline', size: 'lg' }),
                'h-12 px-8 text-base',
              )}
            >
              View Leaderboard
            </Link>
          </div>
        </div>

        <div className="relative z-10 mx-auto mt-12 w-full max-w-6xl px-4 opacity-80 transition-opacity duration-700 hover:opacity-100">
          <div className="overflow-hidden rounded-xl border border-border bg-card/30 shadow-2xl shadow-primary/10 backdrop-blur-sm">
            <HomeChart symbol="bitcoin" />
          </div>
        </div>
      </section>
    </div>
  );
};
