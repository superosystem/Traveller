package com.gusrylmubarok.springbasic;

import com.gusrylmubarok.springbasic.listener.LoginAgainSuccessListener;
import com.gusrylmubarok.springbasic.listener.LoginSuccessListener;
import com.gusrylmubarok.springbasic.listener.UserListener;
import com.gusrylmubarok.springbasic.service.UserService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.context.ConfigurableApplicationContext;
import org.springframework.context.annotation.AnnotationConfigApplicationContext;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Import;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

public class EventListenerTest {

    @Configuration
    @Import({
            UserService.class,
            LoginSuccessListener.class,
            LoginAgainSuccessListener.class,
            UserListener.class
    })
    public static class TestConfiguration {

    }

    private ConfigurableApplicationContext applicationContext;

    @BeforeEach
    void setUp() {
        applicationContext = new AnnotationConfigApplicationContext(TestConfiguration.class);
        applicationContext.registerShutdownHook();
    }

    @Test
    void testEvent() {

        UserService userService = applicationContext.getBean(UserService.class);
        userService.login("eko", "eko");
        userService.login("eko", "salah");
        userService.login("kurninawan", "salah");

    }
}