export const LeaderboardHeader = () => {
  return (
    <thead className="border-b border-border bg-muted text-[10px] font-black uppercase tracking-widest text-muted-foreground">
      <tr>
        <th className="w-24 px-8 py-4 text-center">Rank</th>
        <th className="px-6 py-4">Trader</th>
        <th className="px-8 py-4 text-right">Net Worth (USD)</th>
      </tr>
    </thead>
  );
};
