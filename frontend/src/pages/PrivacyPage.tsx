export const PrivacyPage = () => {
  return (
    <div className="container mx-auto px-4 py-8 max-w-3xl text-foreground pb-24">
      <h1 className="text-3xl font-bold mb-6 text-primary">
        Privacy Policy (Datenschutzerklärung)
      </h1>

      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-2 text-foreground/90">1. Overview</h2>
        <h3 className="text-lg font-medium mb-2 text-foreground/80">General Information</h3>
        <p className="text-muted-foreground mb-4">
          The following notes provide a simple overview of what happens to your personal data when
          you visit this website. Personal data is any data with which you can be personally
          identified.
        </p>
      </section>

      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-2 text-foreground/90">2. Hosting</h2>
        We host the content of our website with the following provider: Google Cloud Platform (GCP).
        The servers are located within the European Union (e.g., Frankfurt, Germany) to ensure
        compliance with GDPR.
      </section>

      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-2 text-foreground/90">
          3. General Notes and Mandatory Information
        </h2>
        <h3 className="text-lg font-medium mb-2 text-foreground/80">Data Protection</h3>
        <p className="text-muted-foreground mb-4">
          The operators of these pages take the protection of your personal data very seriously. We
          treat your personal data confidentially and in accordance with the statutory data
          protection regulations and this data protection declaration.
        </p>
        <h3 className="text-lg font-medium mb-2 text-foreground/80">Responsible Body</h3>
        <p className="text-muted-foreground mb-4">
          The responsible body for data processing on this website is:
          <br />
          <br />
          [TODO: Name]
          <br />
          [TODO: Address, ZIP, City]
          <br />
          E-Mail: [TODO: Email]
        </p>
      </section>

      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-2 text-foreground/90">
          4. Data Collection on this Website
        </h2>
        <h3 className="text-lg font-medium mb-2 text-foreground/80">Cookies</h3>
        <p className="text-muted-foreground mb-4">
          Our website uses so-called "cookies". Cookies are small text files and do not cause any
          damage to your terminal device. They are stored either temporarily for the duration of a
          session (session cookies) or permanently (permanent cookies) on your terminal device.
        </p>
        <p className="text-muted-foreground mb-4">
          <strong>Technically Necessary Cookies (Auth Tokens):</strong> We use cookies that are
          technically necessary for the operation and security of the website (e.g. for
          authentication in the login area using HttpOnly Cookies). These do not require consent.
        </p>

        <h3 className="text-lg font-medium mb-2 text-foreground/80">Server Log Files</h3>
        <p className="text-muted-foreground mb-4">
          The provider of the pages automatically collects and stores information in so-called
          server log files, which your browser automatically transmits to us. These are:
        </p>
        <ul className="list-disc list-inside text-muted-foreground mb-4 ml-4">
          <li>Browser type and browser version</li>
          <li>Operating system used</li>
          <li>Referrer URL</li>
          <li>Hostname of the accessing computer</li>
          <li>Time of the server request</li>
          <li>IP address</li>
        </ul>
        <p className="text-muted-foreground">
          These data are not combined with other data sources.
        </p>
      </section>

      <section className="mb-8">
        <h2 className="text-xl font-semibold mb-2 text-foreground/90">5. User Account</h2>
        <p className="text-muted-foreground mb-4">
          To use functions such as the leaderboard and trading, you must register. We use the data
          entered for the purpose of using the respective offer or service. Registration and usage 
          of the platform require explicit acceptance of our General Terms and Conditions (AGB) and 
          this Privacy Policy.
        </p>
        <p className="text-muted-foreground mb-4">
          <strong>Leaderboards & Public Profiles:</strong> By default, all registered user accounts
          are completely private. While your account might technically rank on our leaderboards based
          on your trading performance, other users cannot view or navigate to your personal profile.
        </p>
        <p className="text-muted-foreground mb-4">
          As an option, you can explicitly choose to make your profile public in your account settings. 
          If you opt-in to a public profile, other users can view your trading portfolio. Additionally, 
          you optionally have the ability to link a personal website or social network to your profile 
          if you wish to share it with others.
        </p>
      </section>
    </div>
  );
};
