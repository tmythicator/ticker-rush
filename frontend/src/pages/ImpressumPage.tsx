export const ImpressumPage = () => {
  return (
    <div className="container mx-auto px-4 py-8 max-w-3xl text-foreground pb-24">
      <h1 className="text-3xl font-bold mb-6 text-primary">Legal Notice (Impressum)</h1>

      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-2 text-foreground/90">
          Information according to § 5 TMG
        </h2>
        <p className="text-muted-foreground">
          {import.meta.env.VITE_LEGAL_NAME || '[VITE_LEGAL_NAME]'}
          <br />
          {import.meta.env.VITE_LEGAL_ADDRESS || '[VITE_LEGAL_ADDRESS]'},{' '}
          {import.meta.env.VITE_LEGAL_CITY || '[VITE_LEGAL_CITY]'}
          <br />
          {import.meta.env.VITE_LEGAL_COUNTRY || '[VITE_LEGAL_COUNTRY]'}
        </p>
      </section>

      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-2 text-foreground/90">Contact</h2>
        <p className="text-muted-foreground">
          E-Mail: {import.meta.env.VITE_LEGAL_EMAIL || '[VITE_LEGAL_EMAIL]'}
        </p>
      </section>

      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-2 text-foreground/90">
          Responsible for Content (V.i.S.d.P.)
        </h2>
        <p className="text-muted-foreground">
          Responsible for content according to § 55 paragraph 2 RStV:
          <br />
          {import.meta.env.VITE_LEGAL_NAME || '[VITE_LEGAL_NAME]'}
          <br />
          {import.meta.env.VITE_LEGAL_ADDRESS || '[VITE_LEGAL_ADDRESS]'},{' '}
          {import.meta.env.VITE_LEGAL_CITY || '[VITE_LEGAL_CITY]'}
          <br />
          {import.meta.env.VITE_LEGAL_COUNTRY || '[VITE_LEGAL_COUNTRY]'}
        </p>
      </section>

      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-2 text-foreground/90">Notice</h2>
        <p className="text-muted-foreground">
          This is a private project for demonstration and educational purposes (Dies ist ein
          privates Projekt zu Demonstrations- und Bildungszwecken). No commercial purpose is pursued
          (Es wird kein wirtschaftlicher Zweck verfolgt).
        </p>
      </section>

      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-2 text-foreground/90">EU Dispute Resolution</h2>
        <p className="text-muted-foreground">
          As this is a private offer and no economic business transaction takes place, we are not
          willing or obliged to participate in dispute resolution proceedings before a consumer
          arbitration board. The EU Online Dispute Resolution (ODR) platform has been discontinued.
        </p>
      </section>
    </div>
  );
};
