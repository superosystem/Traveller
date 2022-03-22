package com.gusrylmubarok.hibernate.util;

import com.gusrylmubarok.hibernate.entity.Name;
import com.gusrylmubarok.hibernate.entity.User;
import org.hibernate.SessionFactory;
import org.hibernate.boot.registry.StandardServiceRegistryBuilder;
import org.hibernate.cfg.Configuration;
import org.hibernate.cfg.Environment;
import org.hibernate.service.ServiceRegistry;

import java.util.Properties;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class HibernateUtil {

    private static SessionFactory sessionFactory;

    public static SessionFactory getSessionFactory() {
        try {
            Configuration configuration = new Configuration();

            // Hibernate setting for connection
            Properties settings = new Properties();
            settings.put(Environment.DRIVER, "org.mysql.cj.jdbc.Driver");
            settings.put(Environment.URL, "jdbc:mysql://localhsot:3306?hibernate_db?useSSL=false&serverTimezone=Asia/Jakarta");
            settings.put(Environment.USER, "root");
            settings.put(Environment.PASS, "");
            settings.put(Environment.DIALECT, "org.hibernate.dialect.MySQL5Dialect");
            settings.put(Environment.SHOW_SQL, "true");
            settings.put(Environment.CURRENT_SESSION_CONTEXT_CLASS, "thread");
            settings.put(Environment.HBM2DDL_AUTO, "create-drop");

            configuration.setProperties(settings);
            configuration.addAnnotatedClass(User.class);

            ServiceRegistry serviceRegistry = new StandardServiceRegistryBuilder()
                    .applySettings(configuration.getProperties()).build();

            sessionFactory = configuration.buildSessionFactory(serviceRegistry);
        } catch (Exception ex) {
            ex.printStackTrace();
        }
        return sessionFactory;
    }

}
