import styles from './Legal.module.css';

export const PrivacyPage = () => {
  return (
    <div className={styles.container}>
      <h1 className={styles.title}>Privacy Policy (Datenschutzerklärung)</h1>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>1. Overview</h2>
        <h3 className={styles.subHeader}>General Information</h3>
        <p className={styles.text}>
          The following notes provide a simple overview of what happens to your personal data when
          you visit this website. Personal data is any data with which you can be personally
          identified.
        </p>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>2. Hosting</h2>
        <p className={styles.text}>
          We host the content of our website with the following provider: Google Cloud Platform
          (GCP). The servers are located within the European Union (e.g., Frankfurt, Germany) to
          ensure compliance with GDPR.
        </p>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>3. General Notes and Mandatory Information</h2>
        <h3 className={styles.subHeader}>Data Protection</h3>
        <p className={styles.text}>
          The operators of these pages take the protection of your personal data very seriously. We
          treat your personal data confidentially and in accordance with the statutory data
          protection regulations and this data protection declaration.
        </p>
        <h3 className={styles.subHeader}>Responsible Body</h3>
        <p className={styles.text}>
          The responsible body for data processing on this website is:
          <br />
          <br />
          {import.meta.env.VITE_LEGAL_NAME || '[VITE_LEGAL_NAME]'}
          <br />
          {import.meta.env.VITE_LEGAL_ADDRESS || '[VITE_LEGAL_ADDRESS]'},{' '}
          {import.meta.env.VITE_LEGAL_CITY || '[VITE_LEGAL_CITY]'}
          <br />
          E-Mail: {import.meta.env.VITE_LEGAL_EMAIL || '[VITE_LEGAL_EMAIL]'}
        </p>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>4. Data Collection on this Website</h2>
        <h3 className={styles.subHeader}>Cookies</h3>
        <p className={styles.text}>
          Our website uses so-called "cookies". Cookies are small text files and do not cause any
          damage to your terminal device. They are stored either temporarily for the duration of a
          session (session cookies) or permanently (permanent cookies) on your terminal device.
        </p>
        <p className={styles.text}>
          <strong>Technically Necessary Cookies (Auth Tokens):</strong> We use cookies that are
          technically necessary for the operation and security of the website (e.g. for
          authentication in the login area using HttpOnly Cookies). These do not require consent.
        </p>

        <h3 className={styles.subHeader}>Server Log Files</h3>
        <p className={styles.text}>
          The provider of the pages automatically collects and stores information in so-called
          server log files, which your browser automatically transmits to us. These are:
        </p>
        <ul className={styles.list}>
          <li>Browser type and browser version</li>
          <li>Operating system used</li>
          <li>Referrer URL</li>
          <li>Hostname of the accessing computer</li>
          <li>Time of the server request</li>
          <li>IP address</li>
        </ul>
        <p className={styles.text}>These data are not combined with other data sources.</p>
      </section>

      <section className={styles.section}>
        <h2 className={styles.sectionHeader}>5. User Account</h2>
        <p className={styles.text}>
          To use functions such as the leaderboard and trading, you must register. We use the data
          entered for the purpose of using the respective offer or service. Registration and usage
          of the platform require explicit acceptance of our General Terms and Conditions (AGB) and
          this Privacy Policy.
        </p>
        <p className={styles.text}>
          <strong>Leaderboards & Public Profiles:</strong> By default, all registered user accounts
          are completely private. While your account might technically rank on our leaderboards
          based on your trading performance, other users cannot view or navigate to your personal
          profile.
        </p>
        <p className={styles.text}>
          As an option, you can explicitly choose to make your profile public in your account
          settings. If you opt-in to a public profile, other users can view your trading portfolio.
          Additionally, you optionally have the ability to link a personal website or social network
          to your profile if you wish to share it with others.
        </p>
      </section>
    </div>
  );
};
