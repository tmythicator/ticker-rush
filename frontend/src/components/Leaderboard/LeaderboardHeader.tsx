export const LeaderboardHeader = () => {
  return (
    <thead className="bg-muted text-muted-foreground uppercase text-[10px] font-black tracking-widest border-b border-border">
      <tr>
        <th className="px-8 py-4 w-24 text-center">Rank</th>
        <th className="px-6 py-4">Trader</th>
        <th className="px-8 py-4 text-right">Net Worth (USD)</th>
      </tr>
    </thead>
  );
};
