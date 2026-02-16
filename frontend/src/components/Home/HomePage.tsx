import { Button } from '@/components/ui/button';
import { IconArrowRight } from '@icons/CustomIcons';
import { Link } from 'react-router-dom';
import { HomeChart } from './HomeChart';

export const HomePage = () => {
  return (
    <div className="flex flex-col min-h-[calc(100vh-4rem)]">
      <section className="flex-1 flex flex-col items-center justify-center relative overflow-hidden pt-20 pb-32">
        <div className="container px-4 md:px-6 relative z-20 flex flex-col items-center text-center">
          <div className="space-y-4 max-w-3xl">
            <h1 className="text-4xl md:text-6xl font-black tracking-tighter sm:text-5xl text-gray-900 dark:bg-clip-text dark:text-transparent dark:bg-gradient-to-r dark:from-white dark:to-gray-500">
              Play the Markets.
              <br />
              <span className="text-primary">Compete in Monthly Ladders.</span>
            </h1>
            <p className="mx-auto max-w-[700px] text-gray-500 md:text-xl dark:text-gray-400">
              Experience the thrill of real-time trading{' '}
              <strong className="font-semibold text-gray-900 dark:text-gray-100">risk-free</strong>.
              Compete in{' '}
              <strong className="font-semibold text-gray-900 dark:text-gray-100">
                monthly ladders
              </strong>{' '}
              and prove your skills to the world.
            </p>
          </div>
          <div className="flex flex-col gap-2 min-[400px]:flex-row mt-8">
            <Button asChild size="lg" className="h-12 px-8 text-base">
              <Link to="/register">
                Get Started <IconArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
            <Button asChild variant="outline" size="lg" className="h-12 px-8 text-base">
              <Link to="/leaderboard">View Leaderboard</Link>
            </Button>
          </div>
        </div>

        <div className="w-full mt-12 relative z-10 max-w-6xl mx-auto px-4 opacity-80 hover:opacity-100 transition-opacity duration-700">
          <div className="rounded-xl overflow-hidden border border-border shadow-2xl shadow-primary/10 bg-card/30 backdrop-blur-sm">
            <HomeChart symbol="CG:bitcoin" />
          </div>
        </div>
      </section>
    </div>
  );
};
