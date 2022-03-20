package com.gusrylmubarok.springbean;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Scope;

/**
 * @author Gusryl Mubarok
 * @email gusrylmubarok@gmail.com
 */

@Configuration
public class AppPrototypeConfiguration {
    @Bean
    @Scope("prototype")
    public UserService userService(){
        return new UserService();
    }
}
